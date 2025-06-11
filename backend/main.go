package main

import (
	"backend/verify"
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"golang.org/x/time/rate"
)

func main() {
	e := echo.New()

	e.Logger.SetLevel(log.WARN)
	e.Logger.SetHeader("${time_rfc3339}")

	e.Use(middleware.BodyLimit("2M"))
	e.Use(middleware.CORS())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${remote_ip} ${method} ${status} ${uri} ${latency_human} ${error}\n",
	}))
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(20))))
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())

	if err := godotenv.Load(); err != nil {
		e.Logger.Fatal(err)
	}

	if os.Getenv("REMMA_NO_VERIFY_REQ") == "yes" {
		verify.REQ = false
	}

	dbPool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer dbPool.Close()

	state := State{DbPool: dbPool}

	e.GET("/metadata", state.GetMetadata)
	e.GET("/relations", state.GetRelations)

	e.GET("/model", state.GetModel)
	e.PUT("/model", state.PutModel)

	e.GET("/search", state.SearchModel)

	if err = e.Start(":8080"); err != nil {
		e.Logger.Fatal(err)
	}
}
