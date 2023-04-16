DB_URL=postgresql://root:pickaxe-db@localhost:5432/pickaxe_db?sslmode=disable

sqlc:
	sqlc generate

db_docs:
	dbdocs build db/docs/pickaxe.dbml

db_schema:
	dbml2sql --postgres -o db/docs/schema.sql db/docs/pickaxe.dbml

postgres:
	docker run --name pickaxe -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=pickaxe-db -d postgres:alpine3.14

docker-network:
	docker network create pickaxe-network

postgres-network:
	docker run --name pickaxe --network pickaxe-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=pickaxe-db -d postgres:alpine3.14

createdb:
	docker exec -it pickaxe createdb --username=root --owner=root pickaxe_db

dropdb:
	docker exec -it pickaxe dropdb pickaxe_db

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

make go:
	go run main.go

build-go:
	go build -o bin/pickaxe -v .

install-go:
	go install

docker-build:
	docker build -t pickaxe:latest .                                                                

docker-container:
	docker run --name pickaxe_app --network pickaxe-network -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE="postgresql://root:pickaxe-db@pickaxe:5432/pickaxe_db?sslmode=disable" pickaxe:latest

.PHONY: sqlc db_docs db_schema postgres docker-network postgres-network createdb migrateup migratedown build-go install-go docker-build docker-container