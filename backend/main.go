package main

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"golang.org/x/time/rate"
)

type Context struct {
	DbPool    *pgxpool.Pool
	Metadata  LiveMetadataConfig
	Relations LiveRelationsConfig
	echo.Context
}

func main() {
	ctx := context.Background()
	e := echo.New()

	err := godotenv.Load()
	if err != nil {
		e.Logger.Fatal(err)
	}

	dbPool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer dbPool.Close()

	cc := Context{DbPool: dbPool}

	e.Logger.SetLevel(log.WARN)
	e.Logger.SetHeader("${time_rfc3339}")

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc.Context = c
			return next(&cc)
		}
	})

	e.Use(middleware.BodyLimit("2M"))
	e.Use(middleware.CORS())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${method} ${status} ${uri} ${error} ${latency_human}\n",
	}))
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(20))))
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())

	e.GET("/metadata", GetMetadata)
	e.GET("/relations", GetRelations)

	e.GET("/model", GetModel)
	e.PUT("/model", PutModel)

	e.GET("/search", SearchModel)

	err = e.Start(":8080")
	if err != nil {
		e.Logger.Fatal(err)
	}
}
