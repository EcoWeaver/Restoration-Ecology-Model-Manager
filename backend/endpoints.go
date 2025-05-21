package main

import (
	"backend/metadata"
	"backend/relations"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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
	err = conn.QueryRow(ctx, "INSERT INTO models (model) $1 RETURNING model_id", graphRequest.Graph, &model_id).Scan(&model_id)
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

	c.Response().Header().Set("Location", fmt.Sprint(model_id))
	return c.NoContent(http.StatusCreated)
}
