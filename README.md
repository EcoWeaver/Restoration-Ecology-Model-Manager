# Restoration Ecology Model Manager

## Running the app

If you are interested in running the app for development purposes, see [CONTRIBUTING.md](CONTRIBUTING.md).

Otherwise, you can follow these steps to use the app locally:

1. Make sure you have [`docker`](https://www.docker.com/) installed. If you don't, you can get it at https://docs.docker.com/get-started/get-docker/.
	- On Windows, `docker` also needs `wsl` to run. To install it, open the Terminal app, type `wsl --update` and press enter.
2. Download the [latest source code](https://github.com/EcoWeaver/Restoration-Ecology-Model-Manager/archive/refs/heads/main.zip), and unzip the downloaded archive.
3. Navigate to the `backend` folder and rename `.env.example` to `.env`.
4. Open a shell in the `Restoration-Ecology-Model-Manager-main` folder.
	- On Windows you can do this from Explorer by holding shift, right-clicking the folder and selecting "Open in Terminal" or "Open PowerShell window here".
	- On MacOS you can do this from Finder by right-clicking the folder and selecting "New Terminal at Folder".
5. In this shell, type `docker compose up --build postgres backend frontend` and press enter.
6. When `docker` has finished building and starting each component, you can access the app at http://localhost:3000.
