package factory

import (
	"context"
	"sync"

	"github.com/sheirys/mine/manager/journal"
	"github.com/sirupsen/logrus"
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

		logrus.WithFields(logrus.Fields{
			"order":  f.Order.ID,
			"action": action,
		}).Info("applying action")

		switch action {
		case ApplyGrinding:
			if err := f.Grind(); err != nil {
				return err
			}
		case ApplySmelting:
			if err := f.Smelt(); err != nil {
				return err
			}
		case ApplyFreezing:
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

	f.listenAndServe()
	return nil

}

func (f *Factory) Stop() {
	f.cancel()
	f.wg.Wait()
	f.conn.Close()
}
