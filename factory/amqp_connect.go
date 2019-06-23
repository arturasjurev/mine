package factory

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

func (f *Factory) listenAMQP() {
	if err := f.prepareRabbit(); err != nil {
		log.Fatal(err)
	}

	go func() {
		f.wg.Add(1)
		for {
			select {
			// handle if message from rabbitmq is received.
			case deliver, _ := <-f.consume:
				task := journal.Order{}
				if err := json.Unmarshal(deliver.Body, &task); err != nil {
					log.Printf("marshal err: %s\n", err)
					continue
				}
				log.Printf("received delivery. id=%s\n", task.ID)
				// factory can process only one order at time and should not
				// queue orders in future. So if HandleOrder(..) returns error
				// that means, this order cannot be processed by this factory at
				// this time. Send NACK to rabbit and requeue order. Maybe other
				// factory is available?
				if err := f.HandleOrder(task); err != nil {
					deliver.Nack(false, true)
					continue
				}
				deliver.Ack(false)

			// handle if application is closed.
			case <-f.ctx.Done():
				f.wg.Done()
				return

			// handle if rabbitmq connection is corrupted. No reconnect logic
			// implemented. If disconnected, kill whole application.
			case <-f.amqpClose:
				f.cancel()
			}
		}
	}()
}

func (f *Factory) prepareRabbit() error {
	var err error
	if f.conn, err = amqp.Dial(f.AMQPAddress); err != nil {
		return fmt.Errorf("failed to connect rabbit: %s", err)
	}
	f.conn.NotifyClose(f.amqpClose)
	if f.ch, err = f.conn.Channel(); err != nil {
		return fmt.Errorf("failed to initiate channel: %s", err)
	}
	if _, err = f.ch.QueueDeclare(ordersAcceptQueue, false, false, false, false, nil); err != nil {
		return fmt.Errorf("failed to declare queue: %s", err)
	}
	if _, err = f.ch.QueueDeclare(ordersStatusQueue, false, false, false, false, nil); err != nil {
		return fmt.Errorf("failed to declare queue: %s", err)
	}
	if f.consume, err = f.ch.Consume(ordersAcceptQueue, "", false, false, false, false, nil); err != nil {
		return fmt.Errorf("failed to consume: %s", err)
	}
	log.Println("connected to rabbitmq")
	return nil
}
