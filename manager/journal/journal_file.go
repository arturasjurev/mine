package journal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// this is journal to file implementation

type JournalFileService struct {
	File         string
	DumpOnChange bool

	data struct {
		Clients []Client `json:"clients"`
		Orders  []Order  `json:"orders"`
	}
}

func (j *JournalFileService) Init() error {
	return nil
}

func (j *JournalFileService) ListClients() ([]Client, error) {
	return j.data.Clients, nil
}

func (j *JournalFileService) Client(ID string) (Client, error) {
	for _, v := range j.data.Clients {
		if v.ID == ID {
			return v, nil
		}
	}
	return Client{}, fmt.Errorf("client %s not found", ID)
}

func (j *JournalFileService) UpsertClient(c Client) (Client, error) {
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

func (j *JournalFileService) ListOrders() ([]Order, error) {
	return j.data.Orders, nil
}

func (j *JournalFileService) Order(ID string) (Order, error) {
	for _, v := range j.data.Orders {
		if v.ID == ID {
			return v, nil
		}
	}
	return Order{}, fmt.Errorf("order %s not found", ID)
}

func (j *JournalFileService) UpsertOrder(o Order) (Order, error) {
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

func (j *JournalFileService) saveToFile() error {
	content, err := json.Marshal(j.data)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(j.File, content, 0644)
}

func (j *JournalFileService) loadFromFile() error {
	content, err := ioutil.ReadFile(j.File)
	if err != nil {
		return err
	}
	return json.Unmarshal(content, &j.data)
}
