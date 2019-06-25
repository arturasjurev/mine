[![Build Status](https://travis-ci.com/sheirys/mine.svg?branch=master)](https://travis-ci.com/sheirys/mine)
[![Go Report Card](https://goreportcard.com/badge/github.com/sheirys/mine)](https://goreportcard.com/report/github.com/sheirys/mine)
[![GoDoc](https://godoc.org/github.com/sheirys/mine?status.svg)](https://godoc.org/github.com/sheirys/mine)
[![codecov](https://codecov.io/gh/sheirys/mine/branch/master/graph/badge.svg)](https://codecov.io/gh/sheirys/mine)

# Mine

This repository tries to satisfy requirements described in [Carbon based life forms](https://github.com/heficed/Carbon-Based-Life-Forms/blob/821ed4bbd7216a8622d6612cad5f50a249ad4f0f/README.md).

* Client is browser. In examples `curl` with `jq` will be used.
* Manager application can be found in `cmd/manager`.
* Factory application can be found in `cmd/factory`.

![mine_datagram](_assets/mine_datagram.svg)

## Manager

Manager stores all client and orders in datafile. This datafile will not
be created on startup and **must be created manually** with:
```bash
    $ echo '{}' > datafile.json
```

By default, manager tries to open `datafile.json` as datafile and starts http
server on `0.0.0.0:8080`. All parameters can be changed via command arguments.
See more with by providing `-h` argument. Manager can be started with:
```bash
    $ go run cmd/manager/main.go
```

Now if manager is started, clients can be created. Possible API endpints can be
found in `manager/http_routes.go` file. New client can be created with:
```bash
    $ curl -s --data-binary '{"name":"client1"}' localhost:8080/clients | jq
    {
    "id": "36ac4626ccb9d5f1",    // client id
    "name": "client1",
    "registered_at": "2019-06-25T23:46:26.966043745+03:00"
    }
```

Now we can create new order for this client. Be sure, that `minetal.state` **match**
`state_from`. Example:
```bash
    $ curl -s --data-binary '{"mineral":{"name":"iron","melting_point":1000,"hardness":2000,"state":"solid"},"state_from":"solid","state_to":"liquid"}' localhost:8080/clients/36ac4626ccb9d5f1/orders | jq
    {
    "id": "d44296a7f6abc6cb",    // order id
    "client_id": "36ac4626ccb9d5f1",
    "finished": false,           // not finished
    "accepted": false,           // not accepted by factory
    "mineral": {
        "name": "iron",
        "state": "solid",
        "melting_point": 1000,
        "hardness": 2000,
        "fractures": 0
    },
    "state_from": "solid",
    "state_to": "liquid",
    "registered_at": "2019-06-25T21:06:30.336031134Z",
    "accepted_at": "0001-01-01T00:00:00Z",
    "finished_at": "0001-01-01T00:00:00Z"
    }
```

Now at any time, this order can be tracked by:
```bash
    $ curl -s localhost:8080/orders/d44296a7f6abc6cb | jq
    {
    "id": "d44296a7f6abc6cb",
    "client_id": "36ac4626ccb9d5f1",
    "finished": false,
    "accepted": false,
    "mineral": {
        "name": "iron",
        "state": "solid",
        "melting_point": 1000,
        "hardness": 2000,
        "fractures": 0
    },
    "state_from": "solid",
    "state_to": "liquid",
    "registered_at": "2019-06-25T21:06:30.336031134Z",
    "accepted_at": "0001-01-01T00:00:00Z",
    "finished_at": "0001-01-01T00:00:00Z"
    }
```

## Factory

Now when we have at least one order probably we want to start factory. Factory
can be started with:
```bash
    $ go run cmd/factory/main.go
```
Factory accepts some cli arguments, be sure to check them with `-h` flag.

When factory is started in console we should see something like:
```bash
$ go run cmd/factory/main.go 
INFO[0000] connected to rabbitmq
INFO[0008] accepted                                      order=d44296a7f6abc6cb
INFO[0008] applying action                               action=apply_grinding order=d44296a7f6abc6cb
INFO[0012] applying action                               action=apply_smelting order=d44296a7f6abc6cb
INFO[0016] finished                                      order=d44296a7f6abc6cb
```
Here, factory accepted order `d44296a7f6abc6cb`. As we created order with
`state_from:"solid"` and `state_to:"liquid"` our mineral transformation recipe
should be `solid->fractured->liquid`. And to reach this state we need to apply
two actions - 1) grinding 2) smelting.

When factory finished, order status can be checked by:
```bash
    $ curl -s localhost:8080/orders/d44296a7f6abc6cb | jq
    {
    "id": "d44296a7f6abc6cb",
    "client_id": "36ac4626ccb9d5f1",
    "finished": true,
    "accepted": true,
    "mineral": {
        "name": "iron",
        "state": "liquid",
        "melting_point": 1000,
        "hardness": 2000,
        "fractures": 0
    },
    "state_from": "solid",
    "state_to": "liquid",
    "registered_at": "2019-06-25T21:06:30.336031134Z",
    "accepted_at": "2019-06-26T00:08:29.47042013+03:00",
    "finished_at": "2019-06-26T00:08:37.471655357+03:00"
    }
```
## Other fun to know

To mimic some realistic factory work you can play by providing different mineral
`hardness` and `melting_points`, or when starting factory provide stronger equipment
like `$ go run cmd/factory/main.go -grinder=100000 -smelter=100000`.

Fast we will realize, that factory is bottleneck if we want to support more
clients. At any time, we can just start more factories. They should scale and
process other queued orders.

## Known issues

* No rabbit reconnection logic implemented.
* Datafile must be created manually if not exist.