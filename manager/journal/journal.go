package journal

import (
	"time"

	"github.com/sheirys/mine/minerals"
)

type Client struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	RegisteredAt time.Time `json:"registered_at"`
}

type Order struct {
	ID           string           `json:"id"`
	ClientID     string           `json:"client_id"`
	Finished     bool             `json:"finished"`
	Accepted     bool             `json:"accepted"`
	Mineral      minerals.Mineral `json:"mineral"`
	StateFrom    minerals.State   `json:"state_from"`
	StateTo      minerals.State   `json:"state_to"`
	RegisteredAt time.Time        `json:"registered_at"`
	AcceptedAt   time.Time        `json:"accepted_at"`
	FinishedAt   time.Time        `json:"finished_at"`
}

type JournalService interface {
	Init() error

	ListClients() ([]Client, error)
	Client(id string) (Client, error)
	UpsertClient(c Client) (Client, error)

	ListOrders() ([]Order, error)
	Order(id string) (Order, error)
	UpsertOrder(o Order) (Order, error)
}
