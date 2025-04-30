package main

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

var postgres_pool *pgxpool.Pool

func main() {
	ctx := context.Background()
	err := error(nil)
	e := echo.New()

	e.GET("/", test)

	postgres_pool, err = pgxpool.New(ctx, "")
	if err != nil {
		e.Logger.Fatal(err)
		return
	}
	defer postgres_pool.Close()

	if err = e.Start(":8080"); err != nil {
		e.Logger.Fatal(err)
		return
	}
}

func test(c echo.Context) error {
	ctx := c.Request().Context()

	conn, err := postgres_pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, "placeholder sql query")
	if err != nil {
		return err
	}
	defer rows.Close()

	return c.String(http.StatusOK, "Hello, World!")
}
