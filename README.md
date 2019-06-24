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
be created on startup and must be created manually with:
```bash
    $ echo '{}' > datafile.json
```

By default, manager tries to open `datafile.json` as datafile and starts http
server on `0.0.0.0:8080`. All parameters can be changed via command arguments.
See more with by providing `-h` argument. Manager can be started with:
```bash
    $ go run cmd/manager/main.go
```

## Factory



## Known issues

* No rabbit reconnection logic implemented.
* Datafile must be created manually if not exist.