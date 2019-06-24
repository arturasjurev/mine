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

	// list all available clients
	// Endpoint: [GET] /clients
	api.Path("/clients").Handler(
		loggerChain.ThenFunc(m.ListClients),
	).Methods(http.MethodGet)

	// create new client
	// Endpoint: [POST] /clients
	api.Path("/clients").Handler(
		loggerChain.ThenFunc(m.CreateClient),
	).Methods(http.MethodPost)

	// list client orders
	// Endpoint: [GET] /clients/{client_id}/orders
	api.Path("/clients/{clientID:[0-9a-f]+}/orders").Handler(
		loggerChain.ThenFunc(m.GetClient),
	).Methods(http.MethodGet)

	// create order
	// Endpoint: [POST] /orders
	api.Path("/orders").Handler(
		loggerChain.ThenFunc(m.CreateOrder),
	).Methods(http.MethodPost)

	// list all orders
	// Endpoint: [GET] /orders
	api.Path("/orders").Handler(
		loggerChain.ThenFunc(m.ListOrders),
	).Methods(http.MethodGet)

	// list single order
	// Endpoint: [GET] /orders/{orderID}
	api.Path("/orders/{orderID:[0-9a-f]+}").Handler(
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
