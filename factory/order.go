package factory

import (
	"log"
	"time"

	"github.com/sheirys/mine/manager/journal"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// ProcessOrder will process order received from manager. Here order states
// accepted/finished will be changed. And manager will be notified about chnaged
// order states.
func (f *Factory) ProcessOrder(o journal.Order, d amqp.Delivery) {

	// factory can process only one order at time and should not queue orders in
	// future. So if factory is in progress now, this order cannot be processed
	// by this factory at this time. Send NACK to rabbit and requeue order.
	// Maybe other factory is available?

	if f.inProgress {
		d.Nack(false, true)
	}

	logrus.WithField("order", o.ID).Info("accepted")

	f.inProgress = true
	f.Order = o

	// notify manager about accepted order
	f.Order.Accepted = true
	f.Order.AcceptedAt = time.Now()
	f.publishState()

	if err := f.Process(); err != nil {
		log.Printf("failed to process: %s\n", err)
	}

	// notify manager about finished order
	f.Order.Finished = true
	f.Order.FinishedAt = time.Now()
	f.publishState()

	f.inProgress = false

	logrus.WithField("order", o.ID).Info("finished")

	// send ack, so rabbitmq will know, that this client can now get other
	// queued order.
	d.Ack(false)

}
