package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

var postgres_pool *pgxpool.Pool

func main() {
	/*
		ctx := context.Background()
		err := error(nil)
		e := echo.New()

		e.GET("/", test)

		// acquire postgres pool
		postgres_pool, err = pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
		if err != nil {
			e.Logger.Fatal(err)
			return
		}
		defer postgres_pool.Close()

		// check credentials are valid and we can connect
		conn, err := postgres_pool.Acquire(ctx)
		if err != nil {
			e.Logger.Fatal(err)
			return
		}
		conn.Release()

		if err = e.Start(":8080"); err != nil {
			e.Logger.Fatal(err)
			return
		}
	*/
	Test_Connection()
}

func test(c echo.Context) error {
	ctx := c.Request().Context()

	conn, err := postgres_pool.Acquire(ctx)
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, "placeholder sql query")
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	defer rows.Close()

	return c.String(http.StatusOK, "Hello, World!")
}

func Test_Connection() {
	// load .env File
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error when loading .env file: %v", err)
	}

	// Get connection sequence from the .env file
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatalf("DATABASE_URL is not correctly placed")
	}

	// Connection to the database is initialized.
	fmt.Println("Connection is initialized...")
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Connection to the database is failed: %v", err)
	}
	defer conn.Close(context.Background())

	// A simple query to test the connection
	var version string
	err = conn.QueryRow(context.Background(), "SELECT version()").Scan(&version)
	if err != nil {
		log.Fatalf("Query can not be processed: %v", err)
	}

	// print database version
	fmt.Println("Successful connection with: ", version)
}
