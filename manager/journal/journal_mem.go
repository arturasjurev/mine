package journal

import (
	"fmt"
)

// MemService implements journal.Journal interface. This implementation uses
// slice in memory to store data about clients and orders.
type MemService struct {
	data struct {
		Clients []Client `json:"clients"`
		Orders  []Order  `json:"orders"`
	}
}

// Init should be called before use of this service. Here various variable
// initializations should be applied.
func (j *MemService) Init() error { return nil }

// ListClients lists all available clients registered in journal. If no
// clients are in journal then empty array should be returned with no error.
func (j *MemService) ListClients() ([]Client, error) {
	return j.data.Clients, nil
}

// Client return single client by privided client id. If no client found, then
// empty client with error will be returned.
func (j *MemService) Client(ID string) (Client, error) {
	for _, v := range j.data.Clients {
		if v.ID == ID {
			return v, nil
		}
	}
	return Client{}, fmt.Errorf("client %s not found", ID)
}

// UpsertClient updates existing client, by client id provided in argument
// client struct. If client is not found, then create new client.
func (j *MemService) UpsertClient(c Client) (Client, error) {
	for i, v := range j.data.Clients {
		if v.ID == c.ID {
			j.data.Clients[i] = c
			return c, nil
		}
	}
	c.ID = generateRandomID()
	j.data.Clients = append(j.data.Clients, c)
	return c, nil
}

// ListOrders lists all available orders registered in journal. If no orders
// exists, then empty array with no error will be returned.
func (j *MemService) ListOrders() ([]Order, error) {
	return j.data.Orders, nil
}

// Order returns single order by provided id. If no order found, empty order
// with error will be returned.
func (j *MemService) Order(ID string) (Order, error) {
	for _, v := range j.data.Orders {
		if v.ID == ID {
			return v, nil
		}
	}
	return Order{}, fmt.Errorf("order %s not found", ID)
}

// UpsertOrder updates existing order by order id provided by argument order.
// If no order found, new order will be created.
func (j *MemService) UpsertOrder(o Order) (Order, error) {
	for i, v := range j.data.Orders {
		if v.ID == o.ID {
			j.data.Orders[i] = o
			return o, nil
		}
	}
	o.ID = generateRandomID()
	j.data.Orders = append(j.data.Orders, o)
	return o, nil
}
