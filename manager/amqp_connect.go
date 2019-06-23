package manager

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/sheirys/mine/manager/journal"
	"github.com/streadway/amqp"
)

const (
	ordersAcceptQueue = "factory-orders-accept"
	ordersStatusQueue = "factory-orders-status"
)

func (m *Manager) listenAMQP() {
	if err := m.prepareRabbit(); err != nil {
		log.Fatal(err)
	}

	go func() {
		m.wg.Add(1)
		for {
			select {
			// handle if order status update received
			case deliver, _ := <-m.consume:
				order := journal.Order{}
				if err := json.Unmarshal(deliver.Body, &order); err != nil {
					log.Printf("marshal err: %s\n", err)
					continue
				}
				m.Journal.UpsertOrder(order)
				log.Printf("order status changed. id=%s accepted=%b finished=%b\n",
					order.ID,
					order.Accepted,
					order.Finished)

			case order := <-m.publish:
				log.Printf("publishing order. id=%s\n", order.ID)
				m.publishOrder(order)

			// handle if application is closed.
			case <-m.ctx.Done():
				m.wg.Done()
				return

			// handle if rabbitmq connection is corrupted. No reconnect logic
			// implemented. If disconnected, kill whole application.
			case <-m.amqpClose:
				m.cancel()
			}
		}
	}()
}

func (m *Manager) prepareRabbit() error {
	var err error
	if m.conn, err = amqp.Dial(m.AMQPAddress); err != nil {
		return fmt.Errorf("failed to connect rabbit: %s", err)
	}
	m.conn.NotifyClose(m.amqpClose)
	if m.ch, err = m.conn.Channel(); err != nil {
		return fmt.Errorf("failed to initiate channel: %s", err)
	}
	if _, err = m.ch.QueueDeclare(ordersAcceptQueue, false, false, false, false, nil); err != nil {
		return fmt.Errorf("failed to declare queue: %s", err)
	}
	if _, err = m.ch.QueueDeclare(ordersStatusQueue, false, false, false, false, nil); err != nil {
		return fmt.Errorf("failed to declare queue: %s", err)
	}
	if m.consume, err = m.ch.Consume(ordersStatusQueue, "", true, false, false, false, nil); err != nil {
		return fmt.Errorf("failed to consume: %s", err)
	}
	return nil
}

func (m *Manager) publishOrder(o journal.Order) error {
	payload, err := json.Marshal(o)
	if err != nil {
		return err
	}
	return m.ch.Publish("",
		ordersAcceptQueue,
		true,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        payload,
		},
	)
}
