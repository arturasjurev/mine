package manager

import "net/http"

func (a *Application) ListClients(w http.ResponseWriter, r *http.Request) {}

func (a *Application) CreateClient(w http.ResponseWriter, r *http.Request) {}

func (a *Application) GetClient(w http.ResponseWriter, r *http.Request) {}

func (a *Application) CreateOrder(w http.ResponseWriter, r *http.Request) {}

func (a *Application) ListOrders(w http.ResponseWriter, r *http.Request) {}

func (a *Application) GetOrder(w http.ResponseWriter, r *http.Request) {}
