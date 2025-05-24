package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

func GetRelations(c echo.Context) error {
	cc := c.(*Context)
	relations, err := cc.Relations.Get()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, relations)
}

func GetMetadata(c echo.Context) error {
	cc := c.(*Context)
	metadata, err := cc.Metadata.Get(c.Logger())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, metadata)
}

func GetModel(c echo.Context) error {
	cc := c.(*Context)
	ctx := c.Request().Context()

	identifier := c.QueryParam("identifier")

	conn, err := cc.DbPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer conn.Rollback(ctx)

	var model Model
	var metadata []Metadata
	err = conn.QueryRow(ctx, "SELECT models.model, jsonb_agg(jsonb_build_object('kind', metadata.kind, 'value', metadata.value)) FROM models INNER JOIN metadata ON models.model_id = metadata.model_id WHERE models.identifier = $1 GROUP BY models.model", identifier).Scan(&model, &metadata)
	if err == pgx.ErrNoRows {
		return c.NoContent(http.StatusNotFound)
	} else if err != nil {
		return err
	}

	err = conn.Commit(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ModelRequest{identifier, model, metadata})
}

func PutModel(c echo.Context) error {
	cc := c.(*Context)
	ctx := c.Request().Context()

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	var modelRequest ModelRequest
	err = json.Unmarshal(body, &modelRequest)
	if err != nil {
		c.Logger().Warn(err)
		return c.NoContent(http.StatusBadRequest)
	}

	err = modelRequest.verify(cc)
	if err != nil {
		c.Logger().Warn(err)
		return c.NoContent(http.StatusBadRequest)
	}

	conn, err := cc.DbPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer conn.Rollback(ctx)

	var modelId int
	err = conn.QueryRow(ctx, "INSERT INTO models (model, identifier) VALUES ($1, $2) RETURNING model_id", modelRequest.Model, modelRequest.Identifier).Scan(&modelId)
	if err != nil {
		return err
	}

	var b pgx.Batch
	for _, metadata := range modelRequest.Metadata {
		b.Queue("INSERT INTO metadata (model_id, kind, value) VALUES ($1, $2, $3)", modelId, metadata.Kind, metadata.Value)
	}
	err = conn.SendBatch(ctx, &b).Close()
	if err != nil {
		return err
	}

	err = conn.Commit(ctx)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)
}

func SearchModel(c echo.Context) error {
	cc := c.(*Context)
	ctx := c.Request().Context()

	query := c.QueryParams()
	if len(query) == 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	var key string
	var val string
	for k, v := range query {
		key = k
		val = v[0]
	}

	conn, err := cc.DbPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer conn.Rollback(ctx)

	var modelId int
	err = conn.QueryRow(ctx, "SELECT model_id FROM metadata WHERE kind = $1 AND value ILIKE $2", key, val).Scan(&modelId)
	if err != nil {
		return err
	}

	var identifier string
	err = conn.QueryRow(ctx, "SELECT identifier FROM models WHERE model_id = $1", modelId).Scan(&identifier)
	if err != nil {
		return err
	}

	err = conn.Commit(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, identifier)
}
