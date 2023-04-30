# pickaxe

Indexer of the Starknet AMM pools written in Golang - to be used in [Fibrous](https://fibrous.finance)

<img src="./pickaxe.png" alt="pickaxe girl" width="250px">

*sister of [Shovel](https://github.com/tahos81/shovel) NFT Indexer*

[![Go Report Card](https://goreportcard.com/badge/github.com/ulerdogan/pickaxe)](https://goreportcard.com/report/github.com/ulerdogan/pickaxe)

<hr/>

The app:

  - follows the Starknet blocks and emits events for the new blocks
  - when new blocks are emitted, it fetches the recorded pools' `sync` events in order to update the pool reserves
  - updates recorded tokens' prices and total pool values periodically
  - will also track the new pools and tokens when the amm contracts become more mature
  
<hr/>

(1) Preperation to build

``` sh
// Prepare the docker network
make docker-network

// Creates Docker container for Postgres
make postgres

// Creates database in the container
make createdb

// Creates Rabbitmq container in Docker
make createdb
```

(2) Build or install the apps: pickaxe & psocket (optional, changes by the running preferences)

``` sh
// Build the apps: psocket & pickaxe
make build

// Install the apps: psocket & pickaxe
make install
```

Run the app directly

``` sh
// Run the socket block finder (basic version - after step-1)
make psocket

// Run the app (basic version - after step-1)
make pickaxe

// Run the app (if the app has been installed - after `install` in step-2)
pickaxe

// Run the socket (if the app has been installed - after `install` step-2)
psocket

// Run the app (if the code has been built - after `build` in step-2)
./bin/pickaxe

// Run the socket (if the code has been built - after `build` in step-2)
./bin/psocket
```


Custom app running preferences 

``` sh
// For testnet envs - run with testnet flag, example:
pickaxe -t
pickaxe --testnet
```

Run the app in docker network (after step-1)

``` sh
// Build the app containers
make docker-build-psocket
make docker-build-pickaxe

// Create & run the app containers
make docker-container-psocket
make docker-container-pickaxe
```

Run the app with docker-compose **(Recommended to run the app - run directly)**


``` sh
// Create the docker compose network
make docker-compose
```

Notes:
* You should prepare a initial amms - tokens - pools list for the initial run. The indexer will accept this point as a synced point. Example ones belo for the mainnet configurations:
  * [amms](./init/states/amms.json)
  * [tokens](./init/states/tokens.json)
  * [pools](./init/states/pools.json)
  * go to [this folder](./init/states_test) to set up testnet configurations

<hr/>

Used Major Requirements & Tools
* [Go](https://go.dev/)
* [dbml](https://dbml-lang.org)
* [docker](https://docker.com/)
* [golang-migrate](https://github.com/golang-migrate/migrate)
* [caigo](https://github.com/dontpanicdao/caigo)
  * [caigo-rpcv02](https://github.com/ulerdogan/caigo-rpcv02) (customized rpcv02 of caigo for pickaxe)
* [gocron](https://github.com/go-co-op/gocron)
* [sqlc](https://sqlc.dev/)
* [rabbitmq](https://www.rabbitmq.com/)

<hr/>

Check the database tables in [DBDocs](https://dbdocs.io/ulerdogan/Pickaxe).

<img src="./db/docs/db_schemas.png" alt="database tables" width="500">