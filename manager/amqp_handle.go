package manager

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/sheirys/mine/manager/journal"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// define rabbitmq queue names.
// FIXME: definitions are duplicated in factory package. Move to separate packet?
// FIXME: see comments on factory/amqp_handle.go
const (
	ordersAcceptQueue = "factory-orders-accept"
	ordersStatusQueue = "factory-orders-status"
)

// listenAndServe will do all rabbitmq related logic for manager. Here we will
// connect to rabbit server, publish new orders to rabbit and consume order
// notifications
func (m *Manager) listenAndServe() {

	if err := m.prepareRabbit(); err != nil {
		log.Fatal(err)
	}

	go func() {
		m.wg.Add(1)
		for {
			select {

			// handle if order status update received
			case d, _ := <-m.consume:
				order := journal.Order{}

				// FIXME: in our scenario order is marshaled with json.Marshal
				// by factory. So maybe we can rely on good marsahler and ignore
				// err here?
				if err := json.Unmarshal(d.Body, &order); err != nil {
					logrus.WithError(err).Error("cannot marshal order update")
					continue
				}

				m.Journal.UpsertOrder(order)

				logrus.WithFields(logrus.Fields{
					"id":       order.ID,
					"accepted": order.Accepted,
					"finished": order.Finished,
				}).Info("order status changed")

			// hande if new order should be published to factory.
			case order := <-m.publish:
				logrus.WithField("id", order.ID).Info("publishing order")
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

// prepareRabbit will prepare rabbit connection to manager.
func (m *Manager) prepareRabbit() error {
	var err error

	// connect to rabbitmq and create channel.
	if m.conn, err = amqp.Dial(m.AMQPAddress); err != nil {
		return fmt.Errorf("failed to connect rabbit: %s", err)
	}
	if m.ch, err = m.conn.Channel(); err != nil {
		return fmt.Errorf("failed to initiate channel: %s", err)
	}

	// create orders accept queue if not created. Manager will put new orders
	// in this queue. Factory will read from this queue and process orders.
	if _, err = m.ch.QueueDeclare(ordersAcceptQueue, false, false, false, false, nil); err != nil {
		return fmt.Errorf("failed to declare queue: %s", err)
	}

	// create orders status queue if not created. Factory will publish
	// notifications about changed order status (e.g. accepted or finished) to
	// this queue, so manager will know when order is finished.
	if _, err = m.ch.QueueDeclare(ordersStatusQueue, false, false, false, false, nil); err != nil {
		return fmt.Errorf("failed to declare queue: %s", err)
	}

	// manager should consume orders from orders status queue.
	if m.consume, err = m.ch.Consume(ordersStatusQueue, "", true, false, false, false, nil); err != nil {
		return fmt.Errorf("failed to consume: %s", err)
	}

	logrus.Info("connected to rabbitmq")
	m.conn.NotifyClose(m.amqpClose)

	return nil
}

// publishOrder will notify factory about new order.
func (m *Manager) publishOrder(o journal.Order) error {

	// here order is created directly from struct, so we can rely
	// on no errors here. Even if error occurs 0f given, because
	// manager should handle invalid JSON.
	payload, _ := json.Marshal(o)
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
