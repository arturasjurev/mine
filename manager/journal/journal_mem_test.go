package journal_test

import (
	"testing"
	"time"

	"github.com/sheirys/mine/manager/journal"
	"github.com/sheirys/mine/minerals"
	"github.com/stretchr/testify/assert"
)

func TestMemServiceInit(t *testing.T) {
	m := &journal.MemService{}
	assert.NoError(t, m.Init())
}

func TestMemServiceClients(t *testing.T) {
	m := &journal.MemService{}
	assert.NoError(t, m.Init())

	create := journal.Client{
		Name:         "some_random_name",
		RegisteredAt: time.Now(),
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

func TestMemServiceNonExistingClient(t *testing.T) {
	m := &journal.MemService{}
	assert.NoError(t, m.Init())

	_, err := m.Client("non_existing_id")
	assert.Error(t, err)
}

func TestMemServiceOrders(t *testing.T) {
	m := &journal.MemService{}
	assert.NoError(t, m.Init())

	create := journal.Client{
		Name:         "some_random_name",
		RegisteredAt: time.Now(),
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

func TestMemServiceNonExistingOrder(t *testing.T) {
	m := &journal.MemService{}
	assert.NoError(t, m.Init())

	_, err := m.Order("non_existing_id")
	assert.Error(t, err)
}

func TestMemServiceListClientOrders(t *testing.T) {
	m := &journal.MemService{}
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
