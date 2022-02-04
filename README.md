# Go-Utils

Go-Utils is a library containing a collection of Golang utilities

## Features

- JSON-RPC client allowing to connect to any JSON-RPC server over HTTP. It is built using [go-autorest](https://github.com/Azure/go-autorest) library, it allows to easily adapt the client to specific server's configuration without having to modify the primary implementation. For example it is possible to add authorization, circuit breakers, request limiters, custom request headers, etc.