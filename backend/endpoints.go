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

	var model_id int
	var model Model
	err = conn.QueryRow(ctx, "SELECT model_id, model FROM models WHERE identifier = $1", identifier).Scan(&model_id, &model)
	if err == pgx.ErrNoRows {
		return c.NoContent(http.StatusNotFound)
	} else if err != nil {
		return err
	}

	rows, err := conn.Query(ctx, "SELECT kind, value FROM metadata_values WHERE model_id = $1", model_id)
	if err != nil {
		return err
	}

	var kind, value string
	var metadata []Metadata
	pgx.ForEachRow(rows, []any{&kind, &value}, func() error {
		metadata = append(metadata, Metadata{kind, value})
		return nil
	})

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

	for _, metadata := range graphRequest.Metadata {
		_, err = conn.Exec(ctx, "INSERT INTO metadata_values (model_id, type, value) VALUES ($1, $2, $3)", model_id, metadata.Kind, metadata.Value)
		if err != nil {
			return err
		}
	}

	err = conn.Commit(ctx)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)
}
