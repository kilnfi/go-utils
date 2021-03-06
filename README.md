# Go-Utils

Go-Utils is a library containing a collection of Golang utilities

## Features

- JSON-RPC client allowing to connect to any JSON-RPC server over HTTP. It is built using [go-autorest](https://github.com/Azure/go-autorest) library, it allows to easily adapt the client to specific server's configuration without having to modify the primary implementation. For example it allows to add authorization, circuit breakers, request limiters, custom request headers, etc.

- Ethereum 1.0 client allowing to connect to any Ethereum node 

    | Features                                                 | Available |
    |----------------------------------------------------------|-----------|
    | Connect to node over HTTP                                | Yes       |
    | Connect to node over WebSocket                           | Not yet   |
    | Use go context for timeout and cancellation              | Yes       |
    | Use core go-ethereum types                               | Yes       |
    | Compatible with abigen generated Smart-Contract bindings | Yes       |
    | Provides tracing for requests                            | Not Yet   |

- Ethereum 2.0 client allowing to connect to any Ethereum 2.0 beacon chain node (compatible with Prysm, Teku, Lighthouse)

    | Features                                                                 | Available |
    |--------------------------------------------------------------------------|-----------|
    | Connect to beacon node over HTTP                                         | Yes       |
    | Use go context for timeout and cancellation                              | Yes       |
    | Use core [protolambda/zrnt](https://github.com/protolambda/zrnt) types   | Yes      |
    | Provides tracing for requests                                            | Not Yet   |

- A collection of Ethereum 1.0 & 2.0 flags compatible with [Cobra](https://github.com/spf13/cobra) library to build CLI applications that need to interact with blockchain nodes

- Helpers to manipulate data into CSV files