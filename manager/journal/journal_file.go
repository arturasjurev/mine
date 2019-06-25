package journal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// FileService implements journal.Journal interface. This implementation used
// json file as database to store information about orders and clients.
// FIXME: DumpOnChange not implemented.
type FileService struct {
	File         string
	DumpOnChange bool

	data struct {
		Clients []Client `json:"clients"`
		Orders  []Order  `json:"orders"`
	}
}

// Init should be called before use of this service. Here various variable
// initializations should be applied.
func (j *FileService) Init() error {
	return j.loadFromFile()
}

// ListClients lists all available clients registered in journal. If no
// clients are in journal then empty array should be returned with no error.
func (j *FileService) ListClients() ([]Client, error) {
	err := j.loadFromFile()
	return j.data.Clients, err
}

// Client return single client by privided client id. If no client found, then
// empty client with error will be returned.
func (j *FileService) Client(ID string) (Client, error) {
	if err := j.loadFromFile(); err != nil {
		return Client{}, err
	}
	for _, v := range j.data.Clients {
		if v.ID == ID {
			return v, nil
		}
	}
	return Client{}, fmt.Errorf("client %s not found", ID)
}

// UpsertClient updates existing client, by client id provided in argument
// client struct. If client is not found, then create new client.
func (j *FileService) UpsertClient(c Client) (Client, error) {
	if err := j.loadFromFile(); err != nil {
		return Client{}, err
	}
	for i, v := range j.data.Clients {
		if v.ID == c.ID {
			j.data.Clients[i] = c
			err := j.saveToFile()
			return c, err
		}
	}
	c.ID = generateRandomID()
	j.data.Clients = append(j.data.Clients, c)
	err := j.saveToFile()
	return c, err
}

// ListOrders lists all available orders registered in journal. If no orders
// exists, then empty array with no error will be returned.
func (j *FileService) ListOrders() ([]Order, error) {
	err := j.loadFromFile()
	return j.data.Orders, err
}

// Order returns single order by provided id. If no order found, empty order
// with error will be returned.
func (j *FileService) Order(ID string) (Order, error) {
	j.loadFromFile()
	for _, v := range j.data.Orders {
		if v.ID == ID {
			return v, nil
		}
	}
	return Order{}, fmt.Errorf("order %s not found", ID)
}

// UpsertOrder updates existing order by order id provided by argument order.
// If no order found, new order will be created.
func (j *FileService) UpsertOrder(o Order) (Order, error) {
	// FIXME: don't like loadFromFile call.
	j.loadFromFile()
	for i, v := range j.data.Orders {
		if v.ID == o.ID {
			j.data.Orders[i] = o
			j.saveToFile()
			return o, nil
		}
	}
	o.ID = generateRandomID()
	j.data.Orders = append(j.data.Orders, o)
	j.saveToFile()
	return o, nil
}

// ListClientOrders list all orders that belong to given client by
// id provided in arguments. If no orders found empty order list with empty
// error will be returned. If client does not exist empty order list
// will be returned with error.
func (j *FileService) ListClientOrders(id string) ([]Order, error) {
	j.loadFromFile()

	orders := []Order{}
	if _, err := j.Client(id); err != nil {
		return orders, err
	}
	for _, v := range j.data.Orders {
		if v.ClientID == id {
			orders = append(orders, v)
		}
	}
	return orders, nil
}

// saveToFile will dump data to file.
func (j *FileService) saveToFile() error {
	content, err := json.Marshal(j.data)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(j.File, content, 0644)
}

// loadFromFile will load data from file.
func (j *FileService) loadFromFile() error {
	content, err := ioutil.ReadFile(j.File)
	if err != nil {
		return err
	}
	if len(content) == 0 {
		content = []byte(`{}`)
	}
	return json.Unmarshal(content, &j.data)
}
