swag = github.com/swaggo/swag/cmd/swag@latest
golangci-lint = github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

_default:
	@echo "usage: make [prepare|build|fmt|lint|backup]"

prepare:
	go mod download
	go install $(swag)

build:
	go run $(swag) init --parseDependency
	go build .

fmt:
	go run $(swag) fmt
	go run $(golangci-lint) fmt

lint:
	go run $(golangci-lint) fmt -d
	go run $(golangci-lint) run

backup:
	mkdir -p postgresql/backup
	docker exec -it postgres pg_dump -U postgres -d postgres --column-inserts > postgresql/backup/$(shell date +%Y-%m-%d_%H-%M-%S).sql
