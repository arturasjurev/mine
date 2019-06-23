package manager

import (
	"net/http"

	"github.com/sheirys/mine/manager/api"
	"github.com/sheirys/mine/manager/journal"
	"github.com/sirupsen/logrus"
)

// list all available clients
// Endpoint: [GET] /clients
func (m *Manager) ListClients(w http.ResponseWriter, r *http.Request) {
	clients, err := m.Journal.ListClients()
	if err != nil {
		api.JSON(w, http.StatusInternalServerError, err)
		return
	}
	api.JSON(w, http.StatusOK, clients)
}

// create new client
// Endpoint: [POST] /clients
func (m *Manager) CreateClient(w http.ResponseWriter, r *http.Request) {
	client := journal.Client{}
	if ok, err := api.BindJSON(r, &client); !ok {
		api.JSON(w, http.StatusBadRequest, nil)
		logrus.WithError(err).Error("bad request")
		return
	}
	created, err := m.Journal.UpsertClient(client)
	if err != nil {
		api.JSON(w, http.StatusInternalServerError, nil)
		logrus.WithError(err).Error("cannot create client")
		return
	}
	api.JSON(w, http.StatusOK, created)
}

// list client orders
// Endpoint: [GET] /clients/{client_id}/orders
func (m *Manager) GetClient(w http.ResponseWriter, r *http.Request) {}

// create order
// Endpoint: [POST] /clients/{client_id}/orders
func (m *Manager) CreateOrder(w http.ResponseWriter, r *http.Request) {
	order := journal.Order{}
	if ok, err := api.BindJSON(r, &order); !ok {
		api.JSON(w, http.StatusBadRequest, nil)
		logrus.WithError(err).Error("bad request")
		return
	}
	created, err := m.Journal.UpsertOrder(order)
	if err != nil {
		api.JSON(w, http.StatusInternalServerError, nil)
		logrus.WithError(err).Error("cannot create order")
		return
	}
	api.JSON(w, http.StatusOK, created)
	m.publish <- created
}

// list all orders
// Endpoint: [GET] /orders
func (m *Manager) ListOrders(w http.ResponseWriter, r *http.Request) {}

// list single order
// Endpoint: [GET] /orders/{orderID}
func (m *Manager) GetOrder(w http.ResponseWriter, r *http.Request) {}
