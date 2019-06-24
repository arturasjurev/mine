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
	return j.loadFromFile()
}

func (j *JournalFileService) ListClients() ([]Client, error) {
	err := j.loadFromFile()
	return j.data.Clients, err
}

func (j *JournalFileService) Client(ID string) (Client, error) {
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

func (j *JournalFileService) UpsertClient(c Client) (Client, error) {
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

func (j *JournalFileService) ListOrders() ([]Order, error) {
	err := j.loadFromFile()
	return j.data.Orders, err
}

func (j *JournalFileService) Order(ID string) (Order, error) {
	j.loadFromFile()
	for _, v := range j.data.Orders {
		if v.ID == ID {
			return v, nil
		}
	}
	return Order{}, fmt.Errorf("order %s not found", ID)
}

func (j *JournalFileService) UpsertOrder(o Order) (Order, error) {
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
