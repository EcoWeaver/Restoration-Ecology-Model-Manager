package main

import (
	"backend/metadata"
	"backend/relations"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

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
		return c.NoContent(http.StatusBadRequest)
	}

	err = graphRequest.validate()
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
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

var (
	ErrEmptyIdentifier          = errors.New("encountered empty identifier")
	ErrNonRelationElement       = errors.New("encountered non-relation element")
	ErrNonMetadataElement       = errors.New("encountered non-metadata element")
	ErrEmptyMetadataValue       = errors.New("encountered empty metadata value")
	ErrMissingMandatoryMetadata = errors.New("missing mandatory metadata")
)

func (self *GraphRequest) validate() error {
	// check identifier isn't empty
	self.Identifier = strings.TrimSpace(self.Identifier)
	if len(self.Identifier) == 0 {
		return ErrEmptyIdentifier
	}

	// check all relation elements are valid
	for _, relation := range self.Graph.Relations {
		if !relation.Kind.IsAElement() {
			return ErrNonRelationElement
		}
	}

	for _, metadata := range self.Metadata {
		// check metadata element is valid
		if !metadata.Kind.IsAElement() {
			return ErrMissingMandatoryMetadata
		}
		// check metadata value isn't empty
		metadata.Value = strings.TrimSpace(metadata.Value)
		if len(metadata.Value) == 0 {
			return ErrEmptyMetadataValue
		}
		if metadata.Ontology != nil {
			ontology := strings.TrimSpace(*metadata.Ontology)
			if len(ontology) == 0 {
				metadata.Ontology = nil
			} else {
				metadata.Ontology = &ontology
			}
		}
	}

	// check all mandatory metadata elements are present
	neededMetadata := map[metadata.Element]struct{}{}
	for _, element := range metadata.ElementValues() {
		if element.Obligation() == metadata.Mandatory {
			neededMetadata[element] = struct{}{}
		}
	}

	for _, metadata := range self.Metadata {
		delete(neededMetadata, metadata.Kind)
	}

	if len(neededMetadata) != 0 {
		return ErrMissingMandatoryMetadata
	}

	return nil
}
