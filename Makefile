.PHONY: up down logs restart
up:
	docker compose up -d
down:
	docker compose down
logs:
	docker compose logs -f app
restart:
	docker compose up --build -d app