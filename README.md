# Restoration Ecology Model Manager App


## Start Program
The three steps have to be performed successively in separate terminals. 

1. Run Docker Container

This project includes a Docker container that runs both the PostgreSQL database and pgAdmin for visually
exploring and managing the database. The Docker container can be executed by the following
command. pdAdmin can be viewed on http://localhost:15432.

```sh
cd backend
make up
```

2. Run Backend

The backend is build with Go, using the Echo web framework to handle RESTful API requests and responses. It can
be executed by the following command and is currently available on http://localhost:8080.

```sh
cd backend
make run
```

3. Run Frontend

The frontend can be executed by the following command and is available on http://localhost:5173/.

```sh
cd frontend
npm install
npm run dev
```

