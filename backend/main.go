package main

import (
	"backend/model"
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
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

	e.GET("/model/relation-keywords", model.RelationKeywords)

	// acquire postgres pool
	var err error
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
