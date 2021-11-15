postgres:
	docker-compose up -d
createdb:
	docker exec -it postgres14 createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres14 dropdb simple_bank
migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down
sqlc:
	sqlc generate
test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown