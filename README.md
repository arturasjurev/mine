[![Build Status](https://travis-ci.com/sheirys/mine.svg?branch=master)](https://travis-ci.com/sheirys/mine)
[![Go Report Card](https://goreportcard.com/badge/github.com/sheirys/mine)](https://goreportcard.com/report/github.com/sheirys/mine)
[![GoDoc](https://godoc.org/github.com/sheirys/mine?status.svg)](https://godoc.org/github.com/sheirys/mine)
[![codecov](https://codecov.io/gh/sheirys/mine/branch/master/graph/badge.svg)](https://codecov.io/gh/sheirys/mine)

# Mine

This repository tries to satisfy requirements described in [Carbon based life forms](https://github.com/heficed/Carbon-Based-Life-Forms/blob/821ed4bbd7216a8622d6612cad5f50a249ad4f0f/README.md).

* Client is browser. In examples `curl` with `jq` will be used.
* Manager application can be found in `cmd/manager`.
* Factory application can be found in `cmd/manager`.

![mine_datagram](_assets/mine_datagram.svg)

## Known issues

* No rabbit reconnetion logic.