.PHONY: build generate migrate run test tidy

generate:
	templ generate

build:
	go build -o app cmd/server/main.go

migrate:
	migrate -path db/migrations -database "sqlite3://db/toy.db" up

run:
	go run cmd/server/main.go

test:
	go test ./...

tidy:
	go mod tidy