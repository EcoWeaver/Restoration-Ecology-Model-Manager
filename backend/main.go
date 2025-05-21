package main

import (
	"backend/metadata"
	"backend/relations"
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

	e.Use(middleware.BodyLimit("2M"))
	e.Use(middleware.CORS())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${method} ${status} ${uri} ${error}\n",
	}))
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(20))))
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())

	e.GET("/", TestConnection)
	e.GET("/metadata-elements", GetMetadataElements)
	e.GET("/model", GetModel)
	e.GET("/relation-elements", GetRelationElements)

	e.PUT("/model", PutModel)

	err := godotenv.Load()
	if err != nil {
		e.Logger.Fatal(err)
	}

	// check every metadata element has an assigned obligation
	for _, element := range metadata.ElementValues() {
		_ = element.Obligation()
	}

	// check every relations element has an assigned obligation
	for _, element := range relations.ElementValues() {
		_ = element.Category()
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
