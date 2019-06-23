package factory

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

const (
	ordersAcceptQueue = "factory-orders-accept"
	ordersStatusQueue = "factory-orders-status"
)

func (f *Factory) listen() {
	if err := f.prepareRabbit(); err != nil {
		log.Fatal(err)
	}

	go func() {
		log.Println("listening routine ... ")
		f.wg.Add(1)
		for {
			select {
			case deliver, ok := <-f.consume:
				if !ok {
					log.Println("disconnected w.consume")
					f.wg.Done()
					return
				}

				log.Printf("received delivery: %s\n", string(deliver.Body))
			case <-f.ctx.Done():
				f.wg.Done()
				return
			}
		}
	}()
}

func (f *Factory) prepareRabbit() error {
	var err error
	if f.conn, err = amqp.Dial(f.AMQPAddress); err != nil {
		return fmt.Errorf("failed to connect rabbit: %s", err)
	}
	if f.ch, err = f.conn.Channel(); err != nil {
		return fmt.Errorf("failed to initiate channel: %s", err)
	}
	if _, err = f.ch.QueueDeclare(ordersAcceptQueue, false, false, false, false, nil); err != nil {
		return fmt.Errorf("failed to declare queue: %s", err)
	}
	if f.consume, err = f.ch.Consume(ordersAcceptQueue, "", true, false, false, false, nil); err != nil {
		return fmt.Errorf("failed to consume: %s", err)
	}
	return nil
}
