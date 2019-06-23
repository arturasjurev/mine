package manager

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func (m *Manager) Routes() *mux.Router {
	api := mux.NewRouter()

	metricsChain := alice.New(m.doMetrics)

	// list all available clients
	// Endpoint: [GET] /clients
	api.Path("/clients").Handler(
		metricsChain.ThenFunc(m.ListClients),
	).Methods(http.MethodGet)

	// create new client
	// Endpoint: [POST] /clients
	api.Path("/clients").Handler(
		metricsChain.ThenFunc(m.CreateClient),
	).Methods(http.MethodPost)

	// list client orders
	// Endpoint: [GET] /clients/{client_id}/orders
	api.Path("/clients/{clientID:[0-9a-f]+}/orders").Handler(
		metricsChain.ThenFunc(m.GetClient),
	).Methods(http.MethodGet)

	// create order
	// Endpoint: [POST] /orders
	api.Path("/orders").Handler(
		metricsChain.ThenFunc(m.CreateOrder),
	).Methods(http.MethodPost)

	// list all orders
	// Endpoint: [GET] /orders
	api.Path("/orders").Handler(
		metricsChain.ThenFunc(m.ListOrders),
	).Methods(http.MethodGet)

	// list single order
	// Endpoint: [GET] /orders/{orderID}
	api.Path("/orders/{orderID:[0-9a-f]+}").Handler(
		metricsChain.ThenFunc(m.GetOrder),
	).Methods(http.MethodGet)

	return api
}
