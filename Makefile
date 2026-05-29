.PHONY: up down logs restart build
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