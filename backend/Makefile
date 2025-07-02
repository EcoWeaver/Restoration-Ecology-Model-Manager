swag = github.com/swaggo/swag/cmd/swag@latest
golangci-lint = github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

_default:
	@echo "usage: make [up|down|psql|fmt|check|build|run|backup]"

up:
	docker compose up

down:
	docker compose down -v

psql:
	docker exec -it postgres psql -U postgres -d postgres

fmt:
	go run $(swag) fmt
	go run $(golangci-lint) fmt

check:
	go run $(golangci-lint) fmt -d
	go run $(golangci-lint) run

prepare:
	go mod download
	go install $(swag)

build:
	go run $(swag) init --parseDependency
	go build .

run:
	go run $(swag) init --parseDependency
	go run .

backup:
	mkdir -p postgresql/backup
	docker exec -it postgres pg_dump -U postgres -d postgres --column-inserts \
		| grep "INSERT INTO " \
		| sed -e "s/INSERT INTO public./INSERT INTO /g" \
		> postgresql/backup/backup_$(shell date +%Y-%m-%d_%H-%M-%S).sql
