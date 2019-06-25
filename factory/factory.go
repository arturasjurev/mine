package factory

import (
	"context"
	"sync"

	"github.com/sheirys/mine/manager/journal"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// Factory is main factory application structure. Factory will accept orders
// sent by manager, process order and notifies manager about changed order
// status (e.g.: accepted / finished).
type Factory struct {

	// Factory equipment. Factory will use this equipment to perform order and
	// manipulate mineral to reach wanted mineral state.
	Grinder Equipment
	Freezer Equipment
	Smelter Equipment

	// Here currently accepted order is stored. Factory can process only one
	// order at time and should not queue orders. When factory accepts order
	// we will set `inProgess` as `true`.
	// FIXME: `inProgress` is not thread safe but in this scale will work.
	Order      journal.Order
	inProgress bool

	// Variables used to connect and handle rabbitmq connection.
	// TODO: reconnect logic not implemented.
	AMQPAddress string
	conn        *amqp.Connection
	ch          *amqp.Channel
	amqpClose   chan *amqp.Error
	consume     <-chan amqp.Delivery

	// Internal factory variables.
	wg     *sync.WaitGroup
	cancel context.CancelFunc
	ctx    context.Context

	// when order changes status to accepted or finished, order should be pushed
	// to this chan. Order pshed to this chan will be sent via rabbitmq to
	// manager as "order status change" notification.
	updates chan journal.Order
}

// Init must be called before starting factory. Various initial setups must be
// done here. E.g. variable initialization.
func (f *Factory) Init() error {
	f.wg = &sync.WaitGroup{}
	f.ctx, f.cancel = context.WithCancel(context.Background())
	f.updates = make(chan journal.Order, 10)
	f.amqpClose = make(chan *amqp.Error)

	return nil
}

// Start rabbitmq connection and start listen incoming orders via rabbitmq.
func (f *Factory) Start() error {
	return f.listenAndServe()
}

// Stop rabbitmq connection and kill factory application.
func (f *Factory) Stop() {
	f.cancel()
	f.wg.Wait()
	f.conn.Close()
}

// Process will be called when new order is accepted from manager. Here we will
// process order, and change mineral state.
func (f *Factory) Process() error {

	// generate required recipe for this order. Here `recipe` is list of actions
	// (order matters) that should be applied to mineral to reach wanted mineral
	// state.
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
			if err := f.grind(); err != nil {
				return err
			}
		case ApplySmelting:
			if err := f.smelt(); err != nil {
				return err
			}
		case ApplyFreezing:
			if err := f.freeze(); err != nil {
				return err
			}
		}
	}

	return nil
}

// grind will be called on `ApplyGrinding` action.
func (f *Factory) grind() (err error) {
	if err = f.Grinder.Insert(f.Order.Mineral); err != nil {
		return err
	}
	if err = f.Grinder.Process(); err != nil {
		return err
	}
	f.Order.Mineral, err = f.Grinder.Takeout()

	return err
}

// freeze will be called on `ApplyFreeze` action.
func (f *Factory) freeze() (err error) {
	if err = f.Freezer.Insert(f.Order.Mineral); err != nil {
		return err
	}
	if err = f.Freezer.Process(); err != nil {
		return err
	}
	f.Order.Mineral, err = f.Freezer.Takeout()

	return err
}

// smelt will be called on `ApplySmelt` action.
func (f *Factory) smelt() (err error) {
	if err = f.Smelter.Insert(f.Order.Mineral); err != nil {
		return err
	}
	if err = f.Smelter.Process(); err != nil {
		return err
	}
	f.Order.Mineral, err = f.Smelter.Takeout()

	return err
}
