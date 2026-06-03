postgres:
	docker run --name postgres18 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=12345 -d postgres:18-alpine

createdb:
	docker exec -it postgres18 createdb --username=root --owner=root go_bank

dropdb:
	docker exec -it postgres18 dropdb go_bank

migrateup:
	migrate -path migrations -database "postgres://root:12345@localhost:5432/go_bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path migrations -database "postgres://root:12345@localhost:5432/go_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path migrations -database "postgres://root:12345@localhost:5432/go_bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path migrations -database "postgres://root:12345@localhost:5432/go_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v ./...

server:
	go run main.go

mock:
	mockgen --destination internal/db/mock/store.go github.com/go-bank/internal/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlc test server mock
