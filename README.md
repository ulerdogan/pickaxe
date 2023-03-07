# pickaxe

Indexer of the Starknet AMM pools written in Golang - to be used in [Fibrous](https://fibrous.finance)

<img src="./pickaxe.png" alt="pickaxe girl" width="250px">

*sister of [Shovel](https://github.com/tahos81/shovel) NFT Indexer*

<hr/>

Preperation to build

``` bash
// Create Docker container for Postgres
make postgres

// Creates database in the container
make createdb

// Install migration tool (macos version)
brew install golang-migrate

// Create database tables
make migrateup
```

<hr/>

Tools
* [dbml](https://dbml-lang.org)
* [docker](https://docker.com/)
* [golang-migrate](https://github.com/golang-migrate/migrate)