postgres:
	docker run --name simplebank -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine
createdb:
	docker exec simplebank createdb --username=root --owner=root bankdb

dropdb:
	docker exec simplebank dropdb bankdb

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/bankdb?sslmode=disable" --verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/bankdb?sslmode=disable" --verbose down

sqlc_init:
	docker run --rm -v $(CURDIR):/src -w /src kjconroy/sqlc init
sqlc_gen:
	docker run --rm -v $(CURDIR):/src -w /src kjconroy/sqlc generate

test:
	go test -v -cover ./...