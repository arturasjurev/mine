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
