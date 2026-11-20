.PHONY: run build migrate-up migrate-down sqlc-gen test

run:
go run cmd/api/main.go

build:
go build -o bin/api cmd/api/main.go

migrate-up:
goose -dir db/migrations postgres "$(DATABASE_URL)" up

migrate-down:
goose -dir db/migrations postgres "$(DATABASE_URL)" down

sqlc-gen:
sqlc generate -f db/sqlc.yaml

test:
go test ./... -v