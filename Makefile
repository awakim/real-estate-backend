postgres:
	docker run --name postgres14 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14.1-alpine3.15

createdb:
	docker exec -it postgres14 createdb --username=root --owner=root immoblock

dropdb:
	docker exec -it postgres14 dropdb immoblock

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/immoblock?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/immoblock?sslmode=disable" -verbose up 1


migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/immoblock?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/immoblock?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/awakim/immoblock-backend/db/sqlc Store

migratecreate:
	migrate create -ext sql -dir db/migration -seq $(migration)

.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 sqlc server mock migratecreate