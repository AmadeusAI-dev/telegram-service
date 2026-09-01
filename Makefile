-include .env
export

up:
	@go run ./cmd/telegram

authorize:
	@go run ./tools/authorize

fmt:
	@go fmt ./...

lint:
	@test -z "$$(gofmt -l .)"

test:
	@go vet ./...
	@go test ./...
