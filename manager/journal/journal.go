package journal

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/sheirys/mine/minerals"
)

// This file contains information about how Journal interface looks like and
// how journal service should work.
//
// Journal is used by manager application as database that stores orders and
// clients.

// Client is someone that can have orders.
type Client struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	RegisteredAt time.Time `json:"registered_at"`
}

// Order can be placed by client to process his mineral. Order will be sent to
// factory. When factory accepts order, it should be marked as `accepted:true`.
// after that, recipe will be generated in factory that will tell, what actions
// should be applied to mineral to reach wanted `state_to` mineral state. When
// all steps from recipe are completed, factory should mark this order as
// finished.
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

// Journal is used by manager to manipulate data with clients and orders.
type Journal interface {

	// Init should be called before use of this service. Here various variable
	// initializations should be applied.
	Init() error

	// ListClients lists all available clients registered in journal. If no
	// clients are in journal then empty array should be returned with no error.
	ListClients() ([]Client, error)

	// Client should return single client by privided client id. If no client
	// found, then empty client with error should be returned.
	Client(id string) (Client, error)

	// UpsertClient should update existing client, by client id provided in
	// argument client. If client is not found, then create new client.
	UpsertClient(c Client) (Client, error)

	// ListOrders should list all available orders registered in journal. If no
	// orders exists, then empty array with no error should be returned.
	ListOrders() ([]Order, error)

	// Order should return single order by provided id. If no order found, empty
	// order with error should be returned.
	Order(id string) (Order, error)

	// UpsertOrder should update existing order by order id provided by argument
	// If order not found, new order should be created.
	UpsertOrder(o Order) (Order, error)
}

// generateRandomID should be used in Journal interface implementations to
// generate order or client id.
func generateRandomID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}
