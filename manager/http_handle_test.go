package manager_test

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sheirys/mine/manager"
	"github.com/sheirys/mine/manager/journal"
	"github.com/stretchr/testify/assert"
)

type mockClient struct {
	api http.Client
}

func (c *mockClient) Do(app *manager.Manager, url, method string, req, data interface{}) (int, error) {
	var (
		reqBody io.Reader
		err     error
		j       []byte
	)

	// Prepare request payload
	if req != nil {
		j, err = json.Marshal(req)
		if err != nil {
			return 0, err
		}
		reqBody = bytes.NewReader(j)
	}

	// Preapre request
	r, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return 0, err
	}

	resp := httptest.NewRecorder()
	app.Routes().ServeHTTP(resp, r)

	// copy response io.ReadCloser for marshaling and verbose
	buf, _ := ioutil.ReadAll(resp.Body)
	body := bytes.NewBuffer(buf)
	if len(buf) <= 0 {
		return resp.Code, nil
	}
	return resp.Code, json.NewDecoder(body).Decode(&data)

}

func testServer() *manager.Manager {
	m := &manager.Manager{
		Journal:       &journal.MemService{},
		DisableRabbit: true,
	}

	m.Init()
	m.Start()

	return m
}

func TestRouteListClients(t *testing.T) {
	client := &mockClient{}
	srv := testServer()
	defer srv.Stop()

	response := []journal.Client{}
	request := journal.Client{
		Name: "some_random_name",
	}

	// new server should have empty clients list
	status, err := client.Do(srv, "/clients", http.MethodGet, nil, &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)
	assert.Len(t, response, 0)

	// create new client
	status, err = client.Do(srv, "/clients", http.MethodPost, request, nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)

	// client list now should contain newly created client
	status, err = client.Do(srv, "/clients", http.MethodGet, nil, &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)
	assert.Len(t, response, 1)

}
