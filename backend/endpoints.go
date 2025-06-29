package main

import (
	"backend/verify"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type State struct {
	DBPool       *pgxpool.Pool
	Metadata     LiveMetadataConfig
	Relations    LiveRelationsConfig
	DOICache     verify.Cache
	SpeciesCache verify.Cache
}

type SearchResult struct {
	UUID       uuid.UUID `json:"uuid"`
	Identifier string    `json:"identifier"`
	Title      string    `json:"title"`
}

// GetMetadata godoc
//
//	@Summary	Get list of currently accepted metadata kinds, types and obligations
//	@Produce	json
//	@Success	200			{object}	MetadataConfig
//	@Failure	500			{object}	echo.HTTPError
//	@Router		/metadata	[get]
func (self *State) GetMetadata(c echo.Context) error {
	metadata, err := self.Metadata.Get(c.Logger())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, metadata)
}

// GetModel godoc
//
//	@Summary		Get a model
//	@Description	Gets a model by its uuid
//	@Produce		json
//	@Param			uuid	query		string	true	"Model UUID"
//	@Success		200		{object}	Model
//	@Failure		404		{object}	echo.HTTPError
//	@Failure		500		{object}	echo.HTTPError
//	@Router			/model	[get]
func (self *State) GetModel(c echo.Context) error {
	ctx := c.Request().Context()

	modelID, err := uuid.Parse(c.QueryParam("uuid"))
	if err != nil {
		c.Logger().Warn(err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	query := `
		SELECT
			models.graph,
			models.parent,
			jsonb_agg(jsonb_build_object('kind', metadata.kind, 'value', metadata.value))
		FROM
			models
			JOIN metadata ON models.model_id = metadata.model_id
		WHERE
			models.model_id = $1
		GROUP BY
			models.graph,
			models.parent`

	conn, err := self.DBPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer conn.Rollback(ctx)

	var graph Graph
	var parent *uuid.UUID
	var metadata []Metadata
	if err = conn.QueryRow(ctx, query, modelID).Scan(&graph, &parent, &metadata); err != nil {
		if err == pgx.ErrNoRows {
			err = fmt.Errorf("model with uuid '%s' not found", modelID)
			c.Logger().Warn(err)
			return c.JSON(http.StatusNotFound, echo.HTTPError{Message: err.Error()})
		} else {
			return err
		}
	}

	if err = conn.Commit(ctx); err != nil {
		return err
	}

	r := Model{&modelID, graph, metadata, parent}

	go func() {
		if err, ok := r.verify(self, c.Logger(), ctx).(verify.Error); ok {
			c.Logger().Warn(fmt.Errorf("encountered error while verifying model with uuid '%s': %w", modelID, err))
		}
	}()

	return c.JSON(http.StatusOK, r)
}

// PutModel godoc
//
//	@Summary		Save a model
//	@Description	Saves a model. The `parent` field in the `Model` is set if it's an edit, null otherwise. The returned model will be identical, except for the `uuid` field, which is updated to the uuid this model was saved as.
//	@Accept			json
//	@Produce		json
//	@Param			model	body		Model	true	"Model"
//	@Success		201		{object}	Model
//	@Failure		400		{object}	echo.HTTPError
//	@Failure		404		{object}	echo.HTTPError
//	@Failure		500		{object}	echo.HTTPError
//	@Router			/model	[put]
func (self *State) PutModel(c echo.Context) error {
	ctx := c.Request().Context()

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	var r Model
	if err = json.Unmarshal(body, &r); err != nil {
		c.Logger().Warn(err)
		return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: err.Error()})
	}

	if err = r.verify(self, c.Logger(), ctx); err != nil {
		if _, ok := err.(verify.Error); ok {
			c.Logger().Warn(err)
			return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: err.Error()})
		} else {
			return err
		}
	}
	conn, err := self.DBPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer conn.Rollback(ctx)

	if r.Parent != nil {
		query := `
		SELECT
			models.graph,
			jsonb_agg(jsonb_build_object('kind', metadata.kind, 'value', metadata.value))
		FROM
			models
			JOIN metadata ON models.model_id = metadata.model_id
		WHERE
			models.model_id = $1
		GROUP BY
			models.graph`

		var graph Graph
		var metadata []Metadata
		if err = conn.QueryRow(ctx, query, r.Parent).Scan(&graph, &metadata); err != nil {
			if err == pgx.ErrNoRows {
				err = fmt.Errorf("encountered non-existent parent model with uuid '%s'", r.Parent)
				c.Logger().Warn(err)
				return c.JSON(http.StatusNotFound, echo.HTTPError{Message: err.Error()})
			} else {
				return err
			}
		}

		if slices.Equal(r.Graph.Nodes, graph.Nodes) && slices.Equal(r.Graph.Relations, graph.Relations) && slices.Equal(r.Metadata, metadata) {
			err = fmt.Errorf("model with uuid '%s' wasn't changed", r.UUID)
			c.Logger().Warn(err)
			return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: err.Error()})
		}
	}

	query := `
		INSERT INTO
			models (graph, parent)
		VALUES
			($1, $2)
		RETURNING
			model_id`

	var modelID uuid.UUID
	if err = conn.QueryRow(ctx, query, r.Graph, r.Parent).Scan(&modelID); err != nil {
		if e, ok := err.(*pgconn.PgError); ok && e.SQLState() == "23503" {
			err = fmt.Errorf("encountered non-existent parent model with uuid '%s'", r.Parent)
			c.Logger().Warn(err)
			return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: err.Error()})
		} else {
			return err
		}
	}
	r.UUID = &modelID

	query = `
		INSERT INTO
			metadata (model_id, kind, value)
		VALUES
			`

	args := []any{modelID}
	var ins []string
	for _, me := range r.Metadata {
		args = append(args, me.Kind, me.Value)
		ins = append(ins, fmt.Sprintf("($1, $%d, $%d)", len(args)-1, len(args)))
	}
	query += strings.Join(ins, `,
			`)

	if _, err = conn.Exec(ctx, query, args...); err != nil {
		return err
	}

	if err = conn.Commit(ctx); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, r)
}

// GetRelations godoc
//
//	@Summary	Get list of currently accepted relation kinds and types
//	@Produce	json
//	@Success	200			{object}	RelationsConfig
//	@Failure	500			{object}	echo.HTTPError
//	@Router		/relations	[get]
func (self *State) GetRelations(c echo.Context) error {
	relations, err := self.Relations.Get(c.Logger())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, relations)
}

// PostSearch godoc
//
//	@Summary		Search for models
//	@Description	Searches for models that match at least one of the provided queries per metadata kind.
//	@Accept			json
//	@Produce		json
//	@Param			queries	body		[]Metadata	true	"Search Queries"
//	@Success		200		{object}	[]SearchResult
//	@Failure		400		{object}	echo.HTTPError
//	@Failure		500		{object}	echo.HTTPError
//	@Router			/search	[post]
func (self *State) PostSearch(c echo.Context) error {
	ctx := c.Request().Context()

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	var r []Metadata
	if err = json.Unmarshal(body, &r); err != nil {
		c.Logger().Warn(err)
		return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: err.Error()})
	}

	metadata, err := self.Metadata.Get(c.Logger())
	if err != nil {
		return err
	}

	params := map[string][]string{}
	for _, q := range r {
		if _, ok := params[q.Kind]; !ok && !metadata.IsMetadataKind(q.Kind) {
			err = fmt.Errorf("encountered invalid metadata kind '%s'", q.Kind)
			c.Logger().Warn(err)
			return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: err.Error()})
		}
		params[q.Kind] = append(params[q.Kind], q.Value)
	}

	query := `
		SELECT
			models.model_id,
			MAX(CASE WHEN metadata.kind = 'identifier' THEN metadata.value END),
			MAX(CASE WHEN metadata.kind = 'title' THEN metadata.value END)
		FROM
			models
			JOIN metadata ON models.model_id = metadata.model_id`

	var args []any
	if len(params) != 0 {
		var conds []string
		for k, v := range params {
			args = append(args, k)
			conj := fmt.Sprintf("(metadata.kind = $%d AND (", len(args))

			var disj []string
			for _, ve := range v {
				args = append(args, "%"+ve+"%")
				disj = append(disj, fmt.Sprintf("metadata.value ILIKE $%d", len(args)))
			}

			conds = append(conds, conj+strings.Join(disj, " OR ")+"))")
		}

		innerQuery := `
				SELECT
					models.model_id
				FROM
					models
					JOIN metadata ON models.model_id = metadata.model_id
				WHERE
					` + strings.Join(conds, `
					OR `)

		if len(params) != 1 {
			args = append(args, len(params))

			innerQuery += fmt.Sprintf(`
				GROUP BY
					models.model_id
				HAVING
					COUNT(DISTINCT metadata.kind) = $%d`, len(args))
		}

		query += `
		WHERE
			models.model_id IN (` + innerQuery + `
			)`
	}

	query += `
		GROUP BY
			models.model_id`

	conn, err := self.DBPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer conn.Rollback(ctx)

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return err
	}

	var modelID uuid.UUID
	var identifier, title string
	var idents []SearchResult
	_, err = pgx.ForEachRow(rows, []any{&modelID, &identifier, &title}, func() error {
		idents = append(idents, SearchResult{modelID, identifier, title})
		return nil
	})
	if err != nil {
		return err
	}

	if err = conn.Commit(ctx); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, idents)
}
