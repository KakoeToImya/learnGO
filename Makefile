.PHONY: up down logs restart build run migrate-create migrate-up migrate-down
include .env
export
DB_URL = postgres://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=disable
up:
	docker compose up -d
down:
	docker compose down
logs:
	docker compose logs -f app
restart:
	docker compose up --build -d app
build:
	docker compose build app
run:
	go run ./cmd/api/main.go
migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name
migrate-up:
	migrate -path migrations -database "$(DB_URL)" up
migrate-down:
	migrate -path migrations -database "$(DB_URL)" down 1
migrate-down-all:
	migrate -path migrations -database "$(DB_URL)" down