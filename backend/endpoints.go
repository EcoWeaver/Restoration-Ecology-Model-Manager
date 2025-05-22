package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

func GetModel(c echo.Context) error {
	ctx := c.Request().Context()

	identifier := c.QueryParam("identifier")

	conn, err := dbPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer conn.Rollback(ctx)

	var model Model
	var metadata []Metadata
	err = conn.QueryRow(ctx, "SELECT m.model, COALESCE(jsonb_agg(jsonb_build_object('kind', mv.kind, 'value', mv.value)), '{}'::jsonb) AS metadata FROM models m LEFT JOIN metadata_values mv ON m.model_id = mv.model_id WHERE m.identifier = $1 GROUP BY m.model", identifier).Scan(&model, &metadata)
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
	ctx := c.Request().Context()

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	var graphRequest ModelRequest
	err = json.Unmarshal(body, &graphRequest)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	conn, err := dbPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer conn.Rollback(ctx)

	var model_id int
	err = conn.QueryRow(ctx, "INSERT INTO models (model, identifier) VALUES ($1, $2) RETURNING model_id", graphRequest.Model, graphRequest.Identifier).Scan(&model_id)
	if err != nil {
		return err
	}

	var b pgx.Batch
	for _, metadata := range graphRequest.Metadata {
		b.Queue("INSERT INTO metadata_values (model_id, kind, value) VALUES ($1, $2, $3)", model_id, metadata.Kind, metadata.Value)
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

	conn, err := dbPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer conn.Rollback(ctx)

	var modelId int
	err = conn.QueryRow(ctx, "SELECT model_id FROM metadata_values WHERE kind = $1 AND value ILIKE $2", key, val).Scan(&modelId)
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
