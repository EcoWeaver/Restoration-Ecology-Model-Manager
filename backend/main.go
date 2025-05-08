package main

import (
	"context"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

var dbPool *pgxpool.Pool

func main() {
	ctx := context.Background()
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.BodyLimit("2M"))
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(20))))

	e.GET("/relation-keywords", RelationKeywords)
	e.GET("/", TestConnection)

	err := godotenv.Load()
	if err != nil {
		e.Logger.Fatal(err)
	}

	// acquire postgres pool
	dbPool, err = pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer dbPool.Close()

	// check credentials are valid and we can connect
	conn, err := dbPool.Acquire(ctx)
	if err != nil {
		e.Logger.Fatal(err)
	}
	conn.Release()

	err = e.Start(":8080")
	if err != nil {
		e.Logger.Fatal(err)
	}
}

func TestConnection(c echo.Context) error {
	ctx := c.Request().Context()

	conn, err := dbPool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	var version string
	err = conn.QueryRow(context.Background(), "SELECT version()").Scan(&version)
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, "Successful connection with: "+version)
}
