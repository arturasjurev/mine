package manager

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/sirupsen/logrus"
)

// Routes binds URL to request handle function.
func (m *Manager) Routes() *mux.Router {
	api := mux.NewRouter()

	loggerChain := alice.New(m.loggerMiddleware)

	// lists all available clients registered in journal.
	// Endpoint: [GET] /clients
	api.Path("/clients").Handler(
		loggerChain.ThenFunc(m.ListClients),
	).Methods(http.MethodGet)

	// creates and registers new client in journal.
	// Endpoint: [POST] /clients
	api.Path("/clients").Handler(
		loggerChain.ThenFunc(m.CreateClient),
	).Methods(http.MethodPost)

	// lists single client from journal by client_id.
	// Endpoint: [GET] /clients/{client_id}/orders
	api.Path("/clients/{client_id:[0-9a-f]+}").Handler(
		loggerChain.ThenFunc(m.GetClient),
	).Methods(http.MethodGet)

	// list all orders from journal that belongs to client.
	// Endpoint: [GET] /clients/{client_id}/orders
	api.Path("/clients/{client_id:[0-9a-f]+}/orders").Handler(
		loggerChain.ThenFunc(m.ListClientOrders),
	).Methods(http.MethodGet)

	// creates and registers new order in journal. And notifies
	// factory about new created order.
	// Endpoint: [POST] /orders
	api.Path("/orders").Handler(
		loggerChain.ThenFunc(m.CreateOrder),
	).Methods(http.MethodPost)

	// lists all orders registered in journal.
	// Endpoint: [GET] /orders
	api.Path("/orders").Handler(
		loggerChain.ThenFunc(m.ListOrders),
	).Methods(http.MethodGet)

	// list single order
	// Endpoint: [GET] /orders/{order_id}
	api.Path("/orders/{order_id:[0-9a-f]+}").Handler(
		loggerChain.ThenFunc(m.GetOrder),
	).Methods(http.MethodGet)

	return api
}

// loggerMiddleware will print in console received HTTP request, method and
// time duration used by called handler.
func (m *Manager) loggerMiddleware(next http.Handler) http.Handler {
	mw := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Now().Sub(start)
		logrus.WithFields(logrus.Fields{
			"method":   r.Method,
			"duration": duration,
		}).Infof("got request %s", r.RequestURI)
	}

	return http.HandlerFunc(mw)
}
