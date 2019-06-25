package manager

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sheirys/mine/manager/api"
	"github.com/sheirys/mine/manager/journal"
	"github.com/sirupsen/logrus"
)

// ListClients lists all available clients registered in journal.
// Endpoint: [GET] /clients
func (m *Manager) ListClients(w http.ResponseWriter, r *http.Request) {
	clients, err := m.Journal.ListClients()
	if err != nil {
		api.JSON(w, http.StatusInternalServerError, err)
		return
	}
	api.JSON(w, http.StatusOK, clients)
}

// CreateClient creates and registers new client in journal.
// Endpoint: [POST] /clients
func (m *Manager) CreateClient(w http.ResponseWriter, r *http.Request) {
	client := journal.Client{}
	if ok, err := api.BindJSON(r, &client); !ok {
		api.JSON(w, http.StatusBadRequest, nil)
		logrus.WithError(err).Error("bad request")
		return
	}
	client.RegisteredAt = time.Now()
	created, err := m.Journal.UpsertClient(client)
	if err != nil {
		api.JSON(w, http.StatusInternalServerError, nil)
		logrus.WithError(err).Error("cannot create client")
		return
	}
	api.JSON(w, http.StatusOK, created)
}

// GetClient lists single client from journal by client_id.
// Endpoint: [GET] /clients/{client_id}
func (m *Manager) GetClient(w http.ResponseWriter, r *http.Request) {
	clients, err := m.Journal.ListClients()
	if err != nil {
		api.JSON(w, http.StatusInternalServerError, err)
		return
	}
	api.JSON(w, http.StatusOK, clients)
}

// ListClientOrders list all orders from journal that belongs to client.
// Endpoint: [GET] /clients/{client_id}/orders
func (m *Manager) ListClientOrders(w http.ResponseWriter, r *http.Request) {
	id := api.SegmentString(mux.Vars(r), "client_id")
	orders, err := m.Journal.ListClientOrders(id)
	if err != nil {
		api.JSON(w, http.StatusNotFound, nil)
		return
	}
	api.JSON(w, http.StatusOK, orders)
}

// CreateOrder creates and registers new order in journal. And notifies factory
// about new created order.
// Endpoint: [POST] /clients/{client_id}/orders
func (m *Manager) CreateOrder(w http.ResponseWriter, r *http.Request) {
	id := api.SegmentString(mux.Vars(r), "client_id")
	order := journal.Order{}
	if ok, err := api.BindJSON(r, &order); !ok {
		api.JSON(w, http.StatusBadRequest, nil)
		logrus.WithError(err).Error("bad request")
		return
	}
	order.ClientID = id
	order.RegisteredAt = time.Now().UTC()
	created, err := m.Journal.UpsertOrder(order)
	if err != nil {
		api.JSON(w, http.StatusInternalServerError, nil)
		logrus.WithError(err).Error("cannot create order")
		return
	}
	api.JSON(w, http.StatusOK, created)
	m.Publish <- created
}

// ListOrders lists all orders registered in journal.
// Endpoint: [GET] /orders
func (m *Manager) ListOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := m.Journal.ListOrders()
	if err != nil {
		api.JSON(w, http.StatusInternalServerError, err)
		return
	}
	api.JSON(w, http.StatusOK, orders)
}

// GetOrder list single order from journal by provided order_id
// Endpoint: [GET] /orders/{order_id}
func (m *Manager) GetOrder(w http.ResponseWriter, r *http.Request) {
	id := api.SegmentString(mux.Vars(r), "order_id")
	order, err := m.Journal.Order(id)
	if err != nil {
		api.JSON(w, http.StatusNotFound, nil)
		return
	}
	api.JSON(w, http.StatusOK, order)
}
