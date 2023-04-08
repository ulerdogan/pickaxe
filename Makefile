DB_URL=postgresql://root:pickaxe-db@localhost:5432/pickaxe_db?sslmode=disable

sqlc:
	sqlc generate

db_docs:
	dbdocs build db/docs/pickaxe.dbml

db_schema:
	dbml2sql --postgres -o db/docs/schema.sql db/docs/pickaxe.dbml

postgres:
	docker run --name pickaxe -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=pickaxe-db -d postgres:alpine3.14

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

.PHONY: sqlc db_docs db_schema postgres createdb migrateup migratedown build-go install-go