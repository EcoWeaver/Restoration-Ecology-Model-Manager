package main

import (
	_ "backend/docs"
	"backend/verify"
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	echoSwagger "github.com/swaggo/echo-swagger/v2"
)

// main godoc
//
//	@title			Restoration Ecology Model Manager API
//	@version		1.0
//	@license.name	AGPL-3.0
//	@license.url	https://www.gnu.org/licenses/agpl-3.0.en.html
//	@host			localhost:8080
func main() {
	e := echo.New()

	e.Use(middleware.BodyLimit(2 * middleware.MB))
	e.Use(middleware.CORS("*"))
	e.Use(middleware.RequestLogger())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())

	if err := godotenv.Load(); err != nil {
		e.Logger.Warn(err.Error())
	}

	noVerifyReq := os.Getenv("REMMA_NO_VERIFY_REQ") == "yes"
	if noVerifyReq {
		e.Logger.Warn("'REMMA_NO_VERIFY_REQ' is set!")
		e.Logger.Warn("This means that metadata verification will not hit external API endpoints!")
		e.Logger.Warn("This must be disabled in a production environment!")
	}

	bioPortalAPIKey := os.Getenv("BIOPORTAL_API_KEY")
	if err := verify.UUID(bioPortalAPIKey); err != nil {
		panic(err)
	}

	dbPool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	defer dbPool.Close()

	state := State{
		ModelMetadata:    HotReload[MetadataConfig]{path: "model_metadata.json"},
		NodeMetadata:     HotReload[MetadataConfig]{path: "node_metadata.json"},
		RelationMetadata: HotReload[MetadataConfig]{path: "relation_metadata.json"},
		Relations:        HotReload[RelationsConfig]{path: "relations.json"},
		DBPool:           dbPool,
		DOICache:         verify.DOICache(noVerifyReq),
		ORCIDCache:       verify.ORCIDCache(noVerifyReq),
		SpeciesCache:     verify.SpeciesCache(noVerifyReq),
		BioPortalAPIKey:  bioPortalAPIKey,
		NoVerifyReq:      noVerifyReq,
	}

	if err := state.runMigrations(context.Background(), e.Logger); err != nil {
		panic(err)
	}

	e.GET("/metadata/model", state.GetMetadataModel)
	e.GET("/metadata/node", state.GetMetadataNode)
	e.GET("/metadata/relation", state.GetMetadataRelation)
	e.GET("/metadata/suggest", state.GetMetadataSuggest)

	e.GET("/model", state.GetModel)
	e.PUT("/model", state.PutModel)
	e.POST("/model/run", state.PostModelRun)
	e.POST("/model/verify", state.PostModelVerify)

	e.GET("/relations", state.GetRelations)

	e.POST("/search", state.PostSearch)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	if err = e.Start(":8080"); err != nil {
		panic(err)
	}
}
