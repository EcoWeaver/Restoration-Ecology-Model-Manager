_default:
	@echo "usage: make [up|down|psql|build|run]"

up:
	docker stop postgres
	docker rm postgres
	docker compose up

down:
	docker compose down

psql:
	psql -h localhost -U postgres -d postgres

build:
	go build .

run:
	go run .
