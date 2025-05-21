package main

import (
	"backend/metadata"
	"backend/relations"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

func GetMetadataElements(c echo.Context) error {
	response := map[metadata.Obligation][]metadata.Element{}

	for _, element := range metadata.ElementValues() {
		obligation := element.Obligation()
		response[obligation] = append(response[obligation], element)
	}

	return c.JSON(http.StatusOK, response)
}

func GetModel(c echo.Context) error {
	ctx := c.Request().Context()

	identifier := c.QueryParam("identifier")

	conn, err := dbPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer conn.Rollback(ctx)

	var graph Graph
	var model_id int
	err = conn.QueryRow(ctx, "SELECT model, model_id FROM models WHERE identifier = $1", identifier).Scan(&graph, &model_id)
	if err == pgx.ErrNoRows {
		return c.NoContent(http.StatusNotFound)
	} else if err != nil {
		return err
	}

	rows, err := conn.Query(ctx, "SELECT element_id, field_value, ontology FROM metadata_values WHERE model_id = $1", model_id)
	if err != nil {
		return err
	}

	var kind metadata.Element
	var value string
	var ontology *string
	var metadata []Metadata
	pgx.ForEachRow(rows, []any{&kind, &value, &ontology}, func() error {
		metadata = append(metadata, Metadata{kind, value, ontology})
		return nil
	})

	err = conn.Commit(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, GraphRequest{identifier, graph, metadata})
}

func GetRelationElements(c echo.Context) error {
	response := map[relations.Category][]relations.Element{}

	for _, element := range relations.ElementValues() {
		category := element.Category()
		response[category] = append(response[category], element)
	}

	return c.JSON(http.StatusOK, response)
}

func PutModel(c echo.Context) error {
	ctx := c.Request().Context()

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	var graphRequest GraphRequest
	err = json.Unmarshal(body, &graphRequest)
	if err != nil {
		return err
	}

	neededMetadata := map[metadata.Element]struct{}{}
	for _, element := range metadata.ElementValues() {
		if element.Obligation() == metadata.Mandatory {
			neededMetadata[element] = struct{}{}
		}
	}

	for _, metadata := range graphRequest.Metadata {
		if !metadata.Kind.IsAElement() {
			return fmt.Errorf("encountered non-metadata element %s", metadata.Kind.String())
		}
		delete(neededMetadata, metadata.Kind)
	}

	if len(neededMetadata) != 0 {
		bytes, err := json.Marshal(neededMetadata)
		if err != nil {
			return err
		}
		return fmt.Errorf("missing mandatory metadata %s", bytes)
	}

	conn, err := dbPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer conn.Rollback(ctx)

	var model_id int
	err = conn.QueryRow(ctx, "INSERT INTO models (model, identifier) VALUES ($1, $2) RETURNING model_id", graphRequest.Graph, graphRequest.Identifier).Scan(&model_id)
	if err != nil {
		return err
	}

	for _, metadata := range graphRequest.Metadata {
		_, err = conn.Exec(ctx, "INSERT INTO metadata_values (model_id, element_id, field_value, ontology) VALUES ($1, $2, $3, $4)", model_id, metadata.Kind, metadata.Value, metadata.Ontology)
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
