# Restoration Ecology Model Manager

## Running the app

If you are interested in running the app for development purposes, see [CONTRIBUTING.md](CONTRIBUTING.md).

Otherwise, you can follow these steps to use the app locally:

1. Make sure you have [`docker`](https://www.docker.com/) installed. If you don't, you can get it at https://docs.docker.com/get-started/get-docker/.
2. Download the [latest source code](https://github.com/EcoWeaver/Restoration-Ecology-Model-Manager/archive/refs/heads/main.zip), and unzip the downloaded archive.
3. Navigate to the `backend` folder and rename `.env.example` to `.env`.
4. Open a terminal in the project root and execute `docker compose up --build postgres backend frontend`.
5. When `docker` has finished building and starting each component, you can access the app at http://localhost:3000.
