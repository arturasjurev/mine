package manager_test

import (
	"io/ioutil"
	"testing"

	"github.com/sheirys/mine/manager"
	"github.com/sheirys/mine/manager/journal"
	"github.com/stretchr/testify/assert"
)

func TestManagerInit(t *testing.T) {
	temp, _ := ioutil.TempFile("", "journal_*")
	m := &manager.Manager{
		Journal: &journal.FileService{
			File: temp.Name(),
		},
		DisableRabbit: true,
	}

	assert.NoError(t, m.Init())
	assert.NoError(t, m.Start())

	m.Stop()
}
