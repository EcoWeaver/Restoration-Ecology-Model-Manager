package main

import (
	"context"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

var postgres_pool *pgxpool.Pool

func main() {
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
