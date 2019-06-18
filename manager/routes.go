package manager

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (a *Application) Routes() *mux.Router {
	api := mux.NewRouter()

	// list all available clients
	// Endpoint: [GET] /clients
	api.Path("/clients").HandlerFunc(a.ListClients).Methods(http.MethodGet)

	// create new client
	// Endpoint: [POST] /clients
	api.Path("/clients").HandlerFunc(a.CreateClient).Methods(http.MethodPost)

	// list client orders
	// Endpoint: [GET] /clients/{client_id}/orders
	api.Path("/clients/{clientID:[0-9a-f]+}/orders").HandlerFunc(a.GetClient).Methods(http.MethodGet)

	// create order
	// Endpoint: [POST] /clients/{client_id}/orders
	api.Path("/clients/{clientID:[0-9a-f]+}/orders").HandlerFunc(a.CreateOrder).Methods(http.MethodPost)

	// list all orders
	// Endpoint: [GET] /orders
	api.Path("/orders").HandlerFunc(a.ListOrders).Methods(http.MethodGet)

	// list single order
	// Endpoint: [GET] /orders/{orderID}
	api.Path("/orders/{orderID:[0-9a-f]+}").HandlerFunc(a.GetOrder).Methods(http.MethodGet)

	return api
}
