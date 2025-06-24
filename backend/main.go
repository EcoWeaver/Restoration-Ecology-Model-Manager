package main

import (
	_ "backend/docs"
	"backend/verify"
	"context"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	echoSwagger "github.com/swaggo/echo-swagger"
	"golang.org/x/time/rate"
)

// @title Restoration Ecology Model Manager API
// @version 1.0

// @license.name AGPL-3.0
// @license.url https://www.gnu.org/licenses/agpl-3.0.en.html

// @host localhost:8080
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

	verify.REQ = os.Getenv("REMMA_NO_VERIFY_REQ") != "yes"

	dbPool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer dbPool.Close()

	state := State{
		DBPool:       dbPool,
		DOICache:     verify.DOICache(),
		ORCIDCache:   verify.ORCIDCache(),
		SpeciesCache: verify.SpeciesCache(),
	}

	e.GET("/metadata", state.GetMetadata)

	e.GET("/model", state.GetModel)
	e.PUT("/model", state.PutModel)

	e.GET("/relations", state.GetRelations)

	e.GET("/search", state.GetSearch)

	e.GET("/swagger", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/versions", state.GetVersions)

	if err = e.Start(":8080"); err != nil {
		e.Logger.Fatal(err)
	}
}
