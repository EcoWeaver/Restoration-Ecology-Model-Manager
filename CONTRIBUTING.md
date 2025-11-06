## Development environment

0. Setup

- To set up the development environment locally, you'll need [`docker`](https://www.docker.com/), [`make`](https://www.gnu.org/software/make/), [`go`](https://go.dev/), [`npm`](https://www.npmjs.com/) and [`nodejs`](https://nodejs.org/) installed.
- The database and backend use an environment file to manage secrets. You'll need to copy the `backend/.env.example` file to `backend/.env`, and uncomment the second `DATABASE_URL` definition for local development.
- You can uncomment the `REMMA_NO_VERIFY_REQ` definition to bypass the backend verification checks that hit an external API endpoint.
- You can find example data for the database at `backend/postgresql/1-example-data.sql.disabled`. To make use of it, copy it to `backend/postgresql/1-example-data.sql` and recreate the `postgres` docker volume.

1. PostgreSQL

- After running the following commands, the PostgreSQL server will be available on http://localhost:5432.
```sh
cd backend
make up
```

2. Backend

- The backend is built with Go, using the Echo web framework. After running the following commands, it will be available on http://localhost:8080.
```sh
cd backend
make run
```

3. Frontend

- The frontend is built with Svelte, using SvelteFlow for node-based user interfaces. After running the following commands, it will be available on http://localhost:5173.
```sh
cd frontend
npm install
npm run dev
```

### dev.sh

If you are using the [`kitty`](https://github.com/kovidgoyal/kitty/) or [`ptyxis`](https://gitlab.gnome.org/chergert/ptyxis/) terminals, or using the [`tmux`](https://github.com/tmux/tmux/) multiplexer, you can use the included `dev.sh` script as a quick-start to replace steps 1-3.

## Linting and formatting

### Backend

The backend uses [`golangci-lint`](https://golangci-lint.run/) to run formatting and linting checks. Formatting can be run with `make fmt` and linting checks can be run with `make lint`.

### Frontend

The frontend uses [`Prettier`](https://prettier.io/) for formatting, [`TypeScript`](https://www.typescriptlang.org/) for type checks and [`ESLint`](https://eslint.org/) for linting checks. Formatting can be run with `npm run format`, type checks can be run with `npm run check` and linting checks can be run with `npm run lint`.

## API documentation

API documentation is built with Swagger, and is available through the frontend by manually navigating to http://localhost:5173/api, or through the backend by manually navigating to http://localhost:8080/swagger/index.html.
