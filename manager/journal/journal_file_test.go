package journal_test

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/sheirys/mine/manager/journal"
	"github.com/sheirys/mine/minerals"
	"github.com/stretchr/testify/assert"
)

// This file contains tests for journal.FileService implementation. As this
// implementation needs file to work properly, on each test case we will
// create temp file a.k.a `mktemp`. Sorry for file flood but I it needs tests.
// After some tests, try to clean `ls -l /tmp/journal_*` files :).

func TestFileServiceInit(t *testing.T) {
	temp, _ := ioutil.TempFile("", "journal_*")
	m := &journal.FileService{
		File: temp.Name(),
	}
	assert.NoError(t, m.Init())
}

func TestFileServiceClients(t *testing.T) {
	temp, _ := ioutil.TempFile("", "journal_*")
	m := &journal.FileService{
		File: temp.Name(),
	}
	assert.NoError(t, m.Init())

	create := journal.Client{
		Name:         "some_random_name",
		RegisteredAt: time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	// create client and check if it will return
	// newly created client.
	created, err := m.UpsertClient(create)
	assert.NoError(t, err)

	assert.NotEqual(t, create.ID, created.ID)
	assert.Equal(t, create.Name, created.Name)
	assert.Equal(t, create.RegisteredAt, created.RegisteredAt)

	// extract created client by id
	got, err := m.Client(created.ID)
	assert.NoError(t, err)
	assert.Equal(t, created, got)

	// update client with new data
	created.Name = "mars_wind"
	updated, err := m.UpsertClient(created)
	assert.NoError(t, err)
	assert.Equal(t, created, updated)

	// check if updated data is returned by id
	got, err = m.Client(created.ID)
	assert.NoError(t, err)
	assert.Equal(t, created, got)

	list, err := m.ListClients()
	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, created, list[0])

}

func TestFileServiceNonExistingClient(t *testing.T) {
	temp, _ := ioutil.TempFile("", "journal_*")
	m := &journal.FileService{
		File: temp.Name(),
	}
	assert.NoError(t, m.Init())

	_, err := m.Client("non_existing_id")
	assert.Error(t, err)
}

func TestFileServiceOrders(t *testing.T) {
	temp, _ := ioutil.TempFile("", "journal_*")
	m := &journal.FileService{
		File: temp.Name(),
	}
	assert.NoError(t, m.Init())

	create := journal.Client{
		Name:         "some_random_name",
		RegisteredAt: time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	created, err := m.UpsertClient(create)

	order := journal.Order{
		ClientID: created.ID,
		Mineral: minerals.Mineral{
			Name:         "jezaus_plaukas",
			State:        minerals.Fracture,
			MeltingPoint: 5000000,
			Hardness:     5000000,
			Fractures:    1,
		},
		StateFrom: minerals.Fracture,
		StateTo:   minerals.Liquid,
	}

	new, err := m.UpsertOrder(order)
	assert.NoError(t, err)

	got, err := m.Order(new.ID)
	assert.NoError(t, err)
	assert.Equal(t, new, got)

	new.Accepted = true
	updated, err := m.UpsertOrder(new)
	assert.NoError(t, err)
	assert.Equal(t, new, updated)

	got, err = m.Order(new.ID)
	assert.NoError(t, err)
	assert.Equal(t, updated, got)

	list, err := m.ListOrders()
	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, updated, list[0])
}

func TestFileServiceNonExistingOrder(t *testing.T) {
	temp, _ := ioutil.TempFile("", "journal_*")
	m := &journal.FileService{
		File: temp.Name(),
	}
	assert.NoError(t, m.Init())

	_, err := m.Order("non_existing_id")
	assert.Error(t, err)
}

func TestFileServiceListClientOrders(t *testing.T) {
	temp, _ := ioutil.TempFile("", "journal_*")
	m := &journal.FileService{
		File: temp.Name(),
	}
	assert.NoError(t, m.Init())

	create := journal.Client{
		Name:         "some_random_name",
		RegisteredAt: time.Now(),
	}

	// create two clients
	client1, err := m.UpsertClient(create)
	assert.NoError(t, err)

	client2, err := m.UpsertClient(create)
	assert.NoError(t, err)

	// create two orders for each client

	order := journal.Order{
		Mineral: minerals.Mineral{
			Name:         "jezaus_plaukas",
			State:        minerals.Fracture,
			MeltingPoint: 5000000,
			Hardness:     5000000,
			Fractures:    1,
		},
		StateFrom: minerals.Fracture,
		StateTo:   minerals.Liquid,
	}

	order.ClientID = client1.ID
	order1, err := m.UpsertOrder(order)
	assert.NoError(t, err)

	order.ClientID = client2.ID
	order2, err := m.UpsertOrder(order)
	assert.NoError(t, err)

	// extract orders for client1
	got, err := m.ListClientOrders(client1.ID)
	assert.NoError(t, err)
	assert.Len(t, got, 1)
	assert.Equal(t, order1, got[0])

	// extract orders for client2
	got, err = m.ListClientOrders(client2.ID)
	assert.NoError(t, err)
	assert.Len(t, got, 1)
	assert.Equal(t, order2, got[0])
}
