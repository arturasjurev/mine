package factory

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/sheirys/mine/manager/journal"
	"github.com/streadway/amqp"
)

type Factory struct {
	// Factory equipment.
	Grinder Equipment
	Freezer Equipment
	Smelter Equipment

	// What order is processed now.
	Order journal.Order

	AMQPAddress string
	conn        *amqp.Connection
	ch          *amqp.Channel
	amqpClose   chan *amqp.Error
	consume     <-chan amqp.Delivery

	wg     *sync.WaitGroup
	cancel context.CancelFunc
	ctx    context.Context

	updates chan journal.Order

	// FIXME: this is not thread safe
	inProgress bool
}

func (f *Factory) Grind() error {
	if err := f.Grinder.Insert(f.Order.Mineral); err != nil {
		return err
	}
	if err := f.Grinder.Process(); err != nil {
		return err
	}
	if product, err := f.Grinder.Takeout(); err != nil {
		return err
	} else {
		f.Order.Mineral = product
	}
	return nil
}

func (f *Factory) Freeze() error {
	if err := f.Freezer.Insert(f.Order.Mineral); err != nil {
		return err
	}
	if err := f.Freezer.Process(); err != nil {
		return err
	}
	if product, err := f.Freezer.Takeout(); err != nil {
		return err
	} else {
		f.Order.Mineral = product
	}
	return nil
}

func (f *Factory) Smelt() error {
	if err := f.Smelter.Insert(f.Order.Mineral); err != nil {
		return err
	}
	if err := f.Smelter.Process(); err != nil {
		return err
	}
	if product, err := f.Smelter.Takeout(); err != nil {
		return err
	} else {
		f.Order.Mineral = product
	}
	return nil
}

func (f *Factory) Process() error {

	recipe, err := GenerateRecipe(f.Order.StateFrom, f.Order.StateTo)
	if err != nil {
		return err
	}
	for _, action := range recipe {
		switch action {
		case ApplyGrinding:
			log.Println("applying grinding")
			if err := f.Grind(); err != nil {
				return err
			}
		case ApplySmelting:
			log.Println("applying smelting")
			if err := f.Smelt(); err != nil {
				return err
			}
		case ApplyFreezing:
			log.Println("applying freezing")
			if err := f.Freeze(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (f *Factory) Init() error {
	f.wg = &sync.WaitGroup{}
	f.ctx, f.cancel = context.WithCancel(context.Background())
	f.updates = make(chan journal.Order, 10)
	f.amqpClose = make(chan *amqp.Error)
	return nil
}

func (f *Factory) Start() error {

	f.listenAMQP()
	return nil

}

func (f *Factory) Stop() {
	f.cancel()
	f.wg.Wait()
	f.conn.Close()
}

func (f *Factory) publishState() error {
	log.Println("publishing state change")
	payload, err := json.Marshal(f.Order)
	if err != nil {
		return err
	}
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
