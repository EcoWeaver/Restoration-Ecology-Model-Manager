# backend/Makefile 

# Start Docker Compose 
up: 
	docker-compose up --build

# Start Docker Container again 
start_again: 
	docker compose up 
# Stop Docker Compose and delete container 
down: 
	docker-compose down -v 

# Connect to database via psql 
psql: 
	psql -h localhost -p 15432 -U user -d postgres_db

# start backend 
run:
	go run main.go