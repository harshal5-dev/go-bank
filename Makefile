postgres:
	docker run --name postgres18 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=12345 -d postgres:18-alpine

createdb:
	docker exec -it postgres18 createdb --username=root --owner=root go_bank

dropdb:
	docker exec -it postgres18 dropdb go_bank

migrateup:
	migrate -path migrations -database "postgres://root:12345@localhost:5432/go_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path migrations -database "postgres://root:12345@localhost:5432/go_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v ./...

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test
