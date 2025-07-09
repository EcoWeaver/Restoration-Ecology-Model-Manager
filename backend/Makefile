swag = github.com/swaggo/swag/cmd/swag@latest
golangci-lint = github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

_default:
	@echo "usage: make [run|build|fmt|prepare|lint|up|down|backup]"

run:
	go run $(swag) init --parseDependency
	go run .

build:
	go run $(swag) init --parseDependency
	go build .

fmt:
	go run $(swag) fmt
	go run $(golangci-lint) fmt

prepare:
	go mod download
	go install $(swag)

lint:
	go run $(golangci-lint) fmt -d
	go run $(golangci-lint) run

up:
	docker compose up postgres

down:
	docker compose down -v postgres

backup:
	mkdir -p postgresql/backup
	docker exec -it postgres pg_dump -U postgres -d postgres --column-inserts > postgresql/backup/$(shell date +%Y-%m-%d_%H-%M-%S).sql
