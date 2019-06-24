package factory

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/sheirys/mine/manager/journal"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// define rabbitmq queue names.
// FIXME: risk with dependency cycle and move this to manager package?
const (
	ordersAcceptQueue = "factory-orders-accept"
	ordersStatusQueue = "factory-orders-status"
)

// listenAndServe will do all rabbitmq related logic for factory. Here we will
// connect to rabbit server, process received orders and notify manager about
// changed order status.
func (f *Factory) listenAndServe() {

	if err := f.prepareRabbit(); err != nil {
		log.Fatal(err)
	}

	go func() {
		f.wg.Add(1)
		for {
			select {

			// handle if new order is received from manager.
			case d, _ := <-f.consume:
				order := journal.Order{}

				// FIXME: in our scenario order is marshaled with json.Marshal
				// by manager. So maybe we can rely on good marsahler and ignore
				// err here?
				if err := json.Unmarshal(d.Body, &order); err != nil {
					log.Printf("marshal err: %s\n", err)
					continue
				}

				// Process received order. We need to pass delivery `d` here
				// because inside ProcessOrder(..) we check if factory is busy
				// so we can send NACK to requeue this order to other factory.
				go f.ProcessOrder(order, d)

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

// prepareRabbit will prepare rabbit connection to factory.
func (f *Factory) prepareRabbit() error {
	var err error

	// connect to rabbitmq and create channel.
	if f.conn, err = amqp.Dial(f.AMQPAddress); err != nil {
		return fmt.Errorf("failed to connect rabbit: %s", err)
	}
	if f.ch, err = f.conn.Channel(); err != nil {
		return fmt.Errorf("failed to initiate channel: %s", err)
	}

	// set fair prefech count to this channel, so rabbitmq will know that this
	// client can process only one delivery at time. In this case, rabbit will
	// not send other deliveries unless ACK or NACK is published back.
	if err := f.ch.Qos(1, 0, false); err != nil {
		return fmt.Errorf("failed to set fair prefetch")
	}

	// create orders accept queue if not created. Manager will put new orders
	// in this queue. Factory will read from this queue and process orders.
	if _, err = f.ch.QueueDeclare(ordersAcceptQueue, false, false, false, false, nil); err != nil {
		return fmt.Errorf("failed to declare queue: %s", err)
	}

	// create orders status queue if not created. Factory will publish
	// notifications about changed order status (e.g. accepted or finished) to
	// this queue, so manager will know when order is finished.
	if _, err = f.ch.QueueDeclare(ordersStatusQueue, false, false, false, false, nil); err != nil {
		return fmt.Errorf("failed to declare queue: %s", err)
	}

	// factory should consume orders from orders accept queue.
	if f.consume, err = f.ch.Consume(ordersAcceptQueue, "", false, false, false, false, nil); err != nil {
		return fmt.Errorf("failed to consume: %s", err)
	}

	logrus.Info("connected to rabbitmq")
	f.conn.NotifyClose(f.amqpClose)

	return nil
}

// publishState will notify manager about changed order status.
func (f *Factory) publishState() error {

	// here order is created directly from struct, so we can rely
	// on no errors here. Even if error occurs 0f given, because
	// manager should handle invalid JSON.
	payload, _ := json.Marshal(f.Order)
	return f.ch.Publish("",
		ordersStatusQueue,
		true,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        payload,
		},
	)
}
