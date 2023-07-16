# go-shelly

Use a subset of the Shelly HTTP API from Golang to control your device.

Currently only features of outlets are supported, but this repository should be fairly easy to extend. It's only HTTP after all!

## How to use?

see [`client_test.go`](client_test.go) :smile:.

## How to build

```
go install github.com/dmarkham/enumer@latest
go generate ./...
go install
```