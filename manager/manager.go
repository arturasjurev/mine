package manager

import (
	"context"
	"net/http"
	"sync"

	"github.com/sheirys/mine/manager/journal"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// Manager is main manager application structure. Manager has HTTP server,
// listens for incoming HTTP API requests, notifies factory about new order and
// collects notifications about processed orders via rabbitmq.
type Manager struct {

	// Journal will store all clients and orders. This is like database for
	// this application. Journal is interfaces, so can be implemented with
	// database, but for this release, journal is simple json file.
	Journal journal.Journal

	// HTTP listen address. E.g. 0.0.0.0:8833
	HTTPAddress string

	// Variables used to connect and handle rabbitmq connection.
	// TODO: reconnect logic not implemented.
	DisableRabbit bool // false for debugging
	AMQPAddress   string
	conn          *amqp.Connection
	ch            *amqp.Channel
	amqpClose     chan *amqp.Error
	consume       <-chan amqp.Delivery

	// Internal manager variables.
	wg     *sync.WaitGroup
	cancel context.CancelFunc
	ctx    context.Context

	// when new order is created order should bu pushed to this chan. Everything
	// pushed to this channel will be published to rabbitmq, so factory should
	// receive this information.
	Publish chan journal.Order
}

// Init must be called before starting manager. Various initial setups must be
// done here. E.g. variable initialization.
func (m *Manager) Init() error {
	m.wg = &sync.WaitGroup{}
	m.ctx, m.cancel = context.WithCancel(context.Background())
	m.Publish = make(chan journal.Order, 5)

	return m.Journal.Init()
}

// Start HTTP server, rabbitmq connection and listen incoming notifications
// about order status changes via rabbitmq.
func (m *Manager) Start() error {

	go m.listenHTTP()
	go m.listenAndServe()

	return nil
}

// Stop rabbitmq connection and kill manager application.
func (m *Manager) Stop() {
	m.cancel()
	m.wg.Wait()
	m.conn.Close()
}

func (m *Manager) listenHTTP() {
	logrus.WithFields(logrus.Fields{
		"addr": m.HTTPAddress,
	}).Info("starting http server")

	http.ListenAndServe(m.HTTPAddress, m.Routes())
}
