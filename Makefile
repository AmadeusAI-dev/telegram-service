include .env
export

up:
	@go run ./cmd/telegram

authorize:
	@go run ./tools/authorize

