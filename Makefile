postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb: 
	docker exec -it postgres12 createdb --username=root --owner=root bank-api

dropdb:
	docker exec -it postgres12 dropdb --username=root --owner=root bank-api

migrateup: 
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/bank-api?sslmode=disable" -verbose up

migrateup1: 
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/bank-api?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/bank-api?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/bank-api?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test: 
	go test -v -cover ./...

server: 
	go run main.go

mock:
	mockgen -destination db/mock/store.go -package mockdb  github.com/noahstew/bank-api/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 sqlc test server mock