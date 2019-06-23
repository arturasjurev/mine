package manager

import (
	"context"
	"net/http"
	"sync"

	"github.com/sheirys/mine/manager/journal"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type Manager struct {
	Journal journal.JournalService

	AMQPAddress string
	conn        *amqp.Connection
	ch          *amqp.Channel
	amqpClose   chan *amqp.Error
	consume     <-chan amqp.Delivery

	wg     *sync.WaitGroup
	cancel context.CancelFunc
	ctx    context.Context

	publish chan journal.Order
}

func (m *Manager) Init() error {
	m.wg = &sync.WaitGroup{}
	m.publish = make(chan journal.Order, 5)
	m.ctx, m.cancel = context.WithCancel(context.Background())

	if err := m.Journal.Init(); err != nil {
		return err
	}
	return nil
}

func (m *Manager) Start() error {

	go m.listenHTTP()
	go m.listenAMQP()

	return nil
}

func (m *Manager) Stop() {

}

func (m *Manager) listenHTTP() {
	logrus.WithFields(logrus.Fields{
		"addr": "0.0.0.0:8833",
	}).Info("starting http server")
	http.ListenAndServe("0.0.0.0:8833", m.Routes())
}
