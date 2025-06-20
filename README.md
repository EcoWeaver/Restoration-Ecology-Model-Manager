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

The backend is built with Go, using the Echo web framework. It can be started with the following commands and is available on http://localhost:8080.

```sh
cd backend
make run
```

3. Run Frontend

The frontend is built with Svelte, using SvelteFlow for node-based user interfaces. It can be started with the following commands and is available on http://localhost:5173/.

```sh
cd frontend
npm install
npm run dev
```
