# Restoration Ecology Model Manager

## Running the app

The three steps have to be performed in separate terminals.

1. Run Docker Container

This project includes a Docker container that runs both the PostgreSQL database and pgAdmin for database administration. The Docker container can be started with the following commands. PostgreSQL is available on http://localhost:25432 and pgAdmin is available on http://localhost:15432.

```sh
cd backend
cp .env.example .env
make up
```

2. Run Backend

The backend is built with Go, using the Echo web framework, and echo-swagger and swag for API documentation. It can be started with the following commands and is available on http://localhost:8080.

```sh
cd backend
make run
```

The documentation can then be viewed at http://localhost:8080/swagger.

3. Run Frontend

The frontend is built with Svelte, using SvelteFlow for node-based user interfaces. It can be started with the following commands and is available on http://localhost:5173/.

```sh
cd frontend
npm install
npm run dev
```

## Linting and formatting

This project uses static analysis tools to ensure correctness and idiomatic code style.

### Backend

The backend uses [`golangci-lint`](https://golangci-lint.run/) to run multiple formatting and linting checks in parallel. Formatting checks can be run with `make fmt` and linting checks can be run with `make check`.

## API documentation

API documentation is built with Swagger, and is available through the frontend by manually navigating to http://localhost:5173/api.
