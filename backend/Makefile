_default:
	@echo "usage: make [up|down|psql|build|run]"

up:
	docker stop postgres || true
	docker rm postgres || true
	docker compose up

down:
	docker compose down

psql:
	docker exec -it postgres psql -U postgres -d postgres

build:
	go build .

run:
	go run .

.PHONY:backup
backup: 
	chmod +x backup.sh
	./backup.sh