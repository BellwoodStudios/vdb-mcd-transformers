# Transformers

## Architecture

Transformers fetch logs from Ethereum, convert/decode them into usable data, and then persist them in postgres.

A transformer consists of:

- A fetcher -> Fetches raw logs from the blockchain and encodes them as go datatypes
- A converter -> Converts this raw data into a human friendly representation suitable for consumption in the API
- A repository -> Abstracts the database

For Maker, vulcanize will be run in `lightSync` mode, so it will store all headers, and then fetchers pull relevant logs by making RPC calls.

## Event Types

For Maker there are two main types of log events that we're tracking:

1. Custom events that are defined in the contract solidity code.
1. `LogNote` events which utilize the [DSNote library](https://github.com/dapphub/ds-note).
1. `Note` events in the `Vat`

The transformer process for each of these different log types is the same, except for the converting process, as denoted below.

## Creating a Transformer

1. Pull an example event (from kovan / ganache etc.)
1. Add event & method sig, contract address, `checked_headers` column name, and label to relevant files in [`constants`](../shared/constants)
1. Write a test for the event sig in [`event_signature_generator_test.go`](../shared/event_signature_generator_test.go)
1. Create DB table (using [`create_migration`](../../scripts/create_migration))
1. Create columns in `checked_headers` in the _same_ migration
1. Add a line to clean the new table `CleanTestDB` (in [`test_config.go`](../../test_config/test_config.go))
1. Define `model.go`
1. Create test event in [`test_data`](../../test_data)
1. Write converter + converter tests
1. Write repository + repository tests
1. Create converter + repository mocks
1. Create an config object [`shared.TransformerConfig`](../shared/transformer.go) in `config.go`
1. Wire up transformer in [`transformers.go`](../transformers.go), remembering to add it to `TransformerInitializers()`
1. Wire up transformer in [`continuousLogSync.go`](../../cmd/continuousLogSync.go)
1. Manually trigger an event and check that it gets persisted to postgres
1. Create an integration test for the shiny new transformer in [`integration_tests`](./integration_tests)

**Fetching Logs**

1. Generate an example raw log event, by either:

   - Pulling the log directly from the Kovan deployment ([address.go](https://github.com/8thlight/maker-vulcanizedb/blob/master/pkg/transformers/shared/constants/address.go)).
   - Deploying the contract to a local chain and emiting the event manually.

1. Fetch the logs from the chain based on the example event's topic zero:

   - The topic zero is based on the keccak-256 hash of the log event's method signature. These are located in [`pkg/transformers/shared/constants/signature.go`](./shared/constants/signature.go).
   - Fetching is done in batch from the [`watcher`](../../libraries/shared/watcher.go).
   - The logs are then chunked up by the [`chunker`](../shared/log_chunker.go) before being delegated to each transformer.

**Coverting logs**

- **Converting most custom events** (such as FlopKick)

  1.  Convert the raw log into a Go struct.
      - We've been using [go-ethereum's abigen tool](https://github.com/ethereum/go-ethereum/tree/master/cmd/abigen) to get the contract's ABI, and a Go struct that represents the event log. We will unpack the raw logs into this struct.
        - To use abigen: `abigen --sol flip.sol --pkg flip --out {/path/to/output_file}`
          - sol: this is the path to the solidity contract
          - pkg: a package name for the generated Go code
          - out: the file path for the generated Go code (optional)
        - the output for `flop.sol` will include the FlopperAbi and the FlopperKick struct:
        ```go
            type FlopperKick struct {
              Id  *big.Int
              Lot *big.Int
              Bid *big.Int
              Gal common.Address
              End *big.Int
              Raw types.Log
            }
        ```
      - Using go-ethereum's `contract.UnpackLog` method we can unpack the raw log into the FlopperKick struct (which we're referring to as the `entity`).
        - See the `ToEntity` method in [`pkg/transformers/flop_kick/converter.go`](./flop_kick/converter.go).
  1.  Convert the entity into a database model. See the `ToModel` method in `pkg/transformers/flop_kick/converter`.

- **Converting Price Feed custom events**

  - Price Feed contracts use the [LogNote event](https://github.com/makerdao/medianizer/blob/master/src/medianizer.sol#L23)
  - The LogNote event takes in the value of the price feed as it's sole argument, and does not index it. This means that this value can be taken directly from the log's data, and then properly converted using the `price_feeds.Convert` method (located in the model.go file).
  - Since this conversion from raw log to model includes less fields than some others, we've chosen to convert it directly to the database model, skipping the `ToEntity` step.

- **Converting LogNote events** (such as tend)
  - Since LogNote events are a generic structure, they depend on the method signature of the method that is calling them. For example, the `tend` method is called on the [flip.sol contract](https://github.com/makerdao/dss/blob/master/src/flip.sol#L117), and it's method signature looks like this: `tend(uint,uint,uint)`.
    - The first four bytes of the Keccak-256 hashed method signature will be located in `topic[0]` on the log.
    - The message sender will be in `topic[1]`.
    - The first parameter passed to `tend` becomes `topic[2]`.
    - The second parameter passed to `tend` will be `topic[3]`.
    - Any additional parameters will be in the log's data field.
    - More detail is located in the [DSNote repo](https://github.com/dapphub/ds-note).

**Get all MissingHeaders**

- Headers are inserted into VulcanizeDB as part of the `lightSync` command. Then for each transformer we check each header for matching logs.
- The MissingHeaders method queries the `checked_headers` table to see if the header has been checked for the given log type.

**Persist the log record to VulcanizeDB**

- Each event log has it's own table in the database, as well as it's own column in the `checked_headers` table.
  - The `checked_headers` table allows us to keep track of which headers have been checked for a given log type.
- To create a new migration file: `./scripts/create_migration create_flop_kick`
  - See `db/migrations/1536942529_create_flop_kick.up.sql`.
  - The specific log event tables are all created in the `maker` schema.
  - There is a one-many association between `headers` and the log
    event tables. This is so that if a header is removed due to a reorg, the associated log event records are also removed.
- To run the migrations: `make migrate HOST=local_host PORT=5432 NAME=vulcanize_private`
- When a new log record is inserted into VulcanizeDB, we also need to make sure to insert a record into the `checked_headers` table for the given log type.
- We have been using the repository pattern (i.e. wrapping all SQL/ORM invocations in isolated namespaces per table) to interact with the database, see the `Create` method in `pkg/transformers/flop_kick/repository.go`.

**MarkHeaderChecked**

- There is a chance that a header does not have a log for the given transformer's log type, and in this instance we also want to record that the header has been "checked" so that we don't continue to query that header over and over.
- In the transformer we'll make sure to insert a row for the header indicating that it has been checked for the log type that the transformer is responsible for.

**Wire each component up in the transformer**

- We use a [`TransformerInitializer`](../shared/transformer.go) struct for each transformer so that we can inject ethRPC and postgresDB connections as well as configuration data (including the contract address, block range, etc.) into the transformer.
- See any of `pkg/transformers/flop_kick/transformer.go`
- All of the transformers are then initialized in `pkg/transformers/transformers.go` with their configuration.
- The transformers can be executed by using the `continuousLogSync` command, which can be configured to run specific transformers or all transformers.

## Useful Documents

[Ethereum Event ABI Specification](https://solidity.readthedocs.io/en/develop/abi-spec.html#events)
