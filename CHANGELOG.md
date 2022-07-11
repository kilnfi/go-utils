# CHANGELOG

### v0.3.0 (July 11th 2022)

- [app] Create app object allowing to orchestrate services

### v0.2.9 (July 7th 2022)

### 🛠️ Bug fixes

- [eth/staking] Fix ValidateDepositData

### v0.2.8 (July 7th 2022)

### 🛠️ Bug fixes

- [eth/staking] Add function to validate DepositData
- [common/types] Fix Marshal Duration

### v0.2.7 (July 4th 2022)

### 🛠️ Bug fixes

- [eth/staking] Add validation on unmarshal of DepositData

### v0.2.6 (July 4th 2022)

### 🛠️ Bug fixes

- [eth/exec] Fix sepolia chain id

### v0.2.5 (July 4th 2022)

### 🛠️ Bug fixes

- [eth/staking] Fix unmarshal/marshal of DepositData

### v0.2.4 (June 28th 2022)

This is an empty release

### v0.2.3 (June 28th 2022)

### :dizzy: Features

- [eth/staking] Add method to verify DepositData

### 🛠️ Bug fixes

- [keystore] Update SignTx to return an error if key is missing

### v0.2.2 (June 23rd 2022)

### :dizzy: Features

- [cmd/keystore] Add command to import keys

### v0.2.1 (June 8th 2022)

### 🕹️ Others

- [mod] Update module name

### v0.2.0 (June 8th 2022)

### :dizzy: Features

- [cmd] Add Cobra commands
- [eth] Refactor ethereum naming from Eth 1 & 2 to execution & consensus layers
- [http] Add various utilities for HTTP
- [hashicorp] Add various utilities for accessing Hashicorp Vault
- [log] Add various utilities aroung logging

### 🛠️ Bug fixes

- [eth2] Fix GetSpec

## v0.1.0 (April 9th 2022)

### :dizzy: Features

- [jsonrpc] Add JSON-RPC client to connect to any JSON-RPC server over HTTP
- [eth1] Add Eth1 client to connect to any Ethereum 1.0 node over HTTP
- [eth2] Add Eth2 client to connect to any Ethereum 2.0 beacon node over HTTP
- [flag] Add a collection of Ethereum 1.0 & 2.0 flags compatible with [Cobra](https://github.com/spf13/cobra) library to build CLI applications
- [csv] Add CSV store to manipulate data in CSV files

### 🕹️ Others

- [ci-cd] Add CI job running unit test
- [ci-cd] Add CI job running lint test