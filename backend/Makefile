_default:
	@echo "usage: make [up|down|psql|check|build|run|backup]"

up:
	docker compose up

down:
	docker compose down

psql:
	docker exec -it postgres psql -U postgres -d postgres

check:
	golangci-lint run
	golangci-lint fmt -d

build:
	go build .

run:
	go run .

backup:
	mkdir -p postgresql/backup
	docker exec -it postgres pg_dump -U postgres -d postgres --column-inserts \
		| grep "INSERT INTO " \
		| sed -e "s/INSERT INTO public./INSERT INTO /g" \
		> postgresql/backup/backup_$(shell date +%Y-%m-%d_%H-%M-%S).sql
