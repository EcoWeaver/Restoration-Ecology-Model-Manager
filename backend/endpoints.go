package main

import (
	"backend/types"
	"backend/util"
	"backend/verify"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type State struct {
	DBPool       *pgxpool.Pool
	Metadata     LiveMetadataConfig
	Relations    LiveRelationsConfig
	DOICache     verify.Cache
	ORCIDCache   verify.Cache
	SpeciesCache verify.Cache
	OLS4Caches   sync.Map
}

type GetModelResult struct {
	UUID     string           `json:"uuid"`
	Uploaded time.Time        `json:"uploaded"`
	Graph    types.Graph      `json:"graph"`
	Parent   *Parent          `json:"parent"`
	Metadata []types.Metadata `json:"metadata"`
}

type PutModelResult struct {
	UUID string `json:"uuid"`
}

type PostModelRunRequest struct {
	UUID            string                 `json:"uuid"`
	NodeValues      map[int]float64        `json:"nodeValues"`
	RelationWeights []types.RelationWeight `json:"relationWeights"`
}

type PostModelRunResult []float64

type PostSearchResult struct {
	UUID       string    `json:"uuid"`
	Uploaded   time.Time `json:"uploaded"`
	Identifier string    `json:"identifier"`
	Title      string    `json:"title"`
	ModelType  string    `json:"modelType"`
}

type Parent struct {
	UUID  string `json:"uuid"`
	Title string `json:"title"`
}

// @Summary		Get the metadata
// @Description	Gets the list of currently accepted metadata kinds, types and obligations.
// @Produce		json
// @Success		200			{object}	MetadataConfig
// @Failure		500			{object}	echo.HTTPError
// @Router		/metadata	[get]
func (self *State) GetMetadata(c echo.Context) error {
	m, err := self.Metadata.Get(c.Logger())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, m)
}

// @Summary		Get a model
// @Description	Gets a model by its UUID.
// @Produce		json
// @Param		uuid	query		string	true	"uuid"
// @Success		200		{object}	GetModelResult
// @Failure		400		{object}	echo.HTTPError
// @Failure		404		{object}	echo.HTTPError
// @Failure		500		{object}	echo.HTTPError
// @Router		/model	[get]
func (self *State) GetModel(c echo.Context) error {
	ctx := c.Request().Context()

	modelID := c.QueryParam("uuid")
	if err := verify.UUID(modelID); err != nil {
		c.Logger().Warn(err)
		return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: err.Error()})
	}

	query := `
		SELECT
			models.uploaded,
			models.graph,
			models.parent,
			jsonb_agg(jsonb_build_object('kind', metadata.kind, 'value', metadata.value))
		FROM
			models
			JOIN metadata ON models.model_id = metadata.model_id
		WHERE
			models.model_id = $1
		GROUP BY
			models.uploaded,
			models.graph,
			models.parent`

	conn, err := self.DBPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer conn.Rollback(ctx)

	var uploaded time.Time
	var graph types.Graph
	var parentUUID *string
	var metadata []types.Metadata
	if err = conn.QueryRow(ctx, query, modelID).Scan(&uploaded, &graph, &parentUUID, &metadata); err != nil {
		if err == pgx.ErrNoRows {
			err = fmt.Errorf("model with uuid '%s' not found", modelID)
			c.Logger().Warn(err)
			return c.JSON(http.StatusNotFound, echo.HTTPError{Message: err.Error()})
		} else {
			return err
		}
	}

	var parent *Parent
	if parentUUID != nil {
		parent = &Parent{UUID: *parentUUID}

		query = `
		SELECT
			metadata.value
		FROM
			metadata
		WHERE
			metadata.model_id = $1 AND metadata.kind = 'title'`

		if err := conn.QueryRow(ctx, query, *parentUUID).Scan(&parent.Title); err != nil {
			return err
		}
	}

	if err = conn.Commit(ctx); err != nil {
		return err
	}

	go func() {
		if err, ok := verifyModel(types.Model{Graph: graph, Parent: parentUUID, Metadata: metadata}, self, c.Logger(), ctx).(types.Error); ok {
			c.Logger().Warn(fmt.Errorf("encountered error while verifying model with uuid '%s': %w", modelID, err))
		}
	}()

	return c.JSON(http.StatusOK, GetModelResult{modelID, uploaded, graph, parent, metadata})
}

// @Summary		Save a model
// @Description	Saves a model. Returns the UUID assigned to that model.
// @Accept		json
// @Produce		json
// @Param		model	body		types.Model	true	"model"
// @Success		201		{object}	PutModelResult
// @Failure		400		{object}	echo.HTTPError
// @Failure		409		{object}	echo.HTTPError
// @Failure		500		{object}	echo.HTTPError
// @Router		/model	[put]
func (self *State) PutModel(c echo.Context) error {
	ctx := c.Request().Context()

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	var m types.Model
	if err = json.Unmarshal(body, &m); err != nil {
		c.Logger().Warn(err)
		return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: err.Error()})
	}

	if err = verifyModel(m, self, c.Logger(), ctx); err != nil {
		if _, ok := err.(types.Error); ok {
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

	if m.Parent != nil {
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

		var graph types.Graph
		var metadata []types.Metadata
		if err := conn.QueryRow(ctx, query, m.Parent).Scan(&graph, &metadata); err != nil {
			if err == pgx.ErrNoRows {
				err = fmt.Errorf("encountered non-existent parent model with uuid '%s'", *m.Parent)
				return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: err.Error()})
			} else {
				return err
			}
		}

		if slices.Equal(m.Graph.Nodes, graph.Nodes) && util.SlicesEqualUnordered(m.Graph.Relations, graph.Relations) && util.SlicesEqualUnordered(m.Metadata, metadata) {
			err = fmt.Errorf("model with uuid '%s' wasn't changed", *m.Parent)
			c.Logger().Warn(err)
			return c.JSON(http.StatusConflict, echo.HTTPError{Message: err.Error()})
		}
	}

	query := `
		INSERT INTO
			models (graph, parent)
		VALUES
			($1, $2)
		RETURNING
			model_id`

	var modelID string
	if err = conn.QueryRow(ctx, query, m.Graph, m.Parent).Scan(&modelID); err != nil {
		return err
	}

	query = `
		INSERT INTO
			metadata (model_id, kind, value)
		VALUES
			`

	args := []any{modelID}
	var vals []string
	for _, me := range m.Metadata {
		args = append(args, me.Kind, me.Value)
		vals = append(vals, fmt.Sprintf("($1, $%d, $%d)", len(args)-1, len(args)))
	}
	query += strings.Join(vals, `,
			`)

	if _, err = conn.Exec(ctx, query, args...); err != nil {
		return err
	}

	if err = conn.Commit(ctx); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, PutModelResult{modelID})
}

// @Summary		Get the relations
// @Description	Gets the list of currently accepted relation kinds and types.
// @Produce		json
// @Success		200			{object}	RelationsConfig
// @Failure		500			{object}	echo.HTTPError
// @Router		/relations	[get]
func (self *State) GetRelations(c echo.Context) error {
	r, err := self.Relations.Get(c.Logger())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, r)
}

// @Summary		Search for models
// @Description	Searches for models that match at least one of the provided queries per metadata kind.
// @Accept		json
// @Produce		json
// @Param		queries	body		[]types.Metadata	true	"queries"
// @Success		200		{object}	[]PostSearchResult
// @Failure		400		{object}	echo.HTTPError
// @Failure		500		{object}	echo.HTTPError
// @Router		/search	[post]
func (self *State) PostSearch(c echo.Context) error {
	ctx := c.Request().Context()

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	var mqs []types.Metadata
	if err = json.Unmarshal(body, &mqs); err != nil {
		c.Logger().Warn(err)
		return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: err.Error()})
	}

	m, err := self.Metadata.Get(c.Logger())
	if err != nil {
		return err
	}

	qs := map[string][]string{}
	for _, mq := range mqs {
		if _, ok := qs[mq.Kind]; !ok && !m.IsMetadataKind(mq.Kind) {
			err = fmt.Errorf("encountered invalid metadata kind '%s'", mq.Kind)
			c.Logger().Warn(err)
			return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: err.Error()})
		}
		qs[mq.Kind] = append(qs[mq.Kind], mq.Value)
	}

	query := `
		SELECT
			models.model_id,
			models.uploaded,
			MAX(CASE WHEN metadata.kind = 'identifier' THEN metadata.value END),
			MAX(CASE WHEN metadata.kind = 'title' THEN metadata.value END),
			MAX(CASE WHEN metadata.kind = 'model type' THEN metadata.value END)
		FROM
			models
			JOIN metadata ON models.model_id = metadata.model_id`

	var args []any
	if len(qs) != 0 {
		var conds []string
		for k, v := range qs {
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

		if len(qs) != 1 {
			args = append(args, len(qs))

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
			models.model_id,
			models.uploaded`

	conn, err := self.DBPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer conn.Rollback(ctx)

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return err
	}

	var sr PostSearchResult
	var res []PostSearchResult
	_, err = pgx.ForEachRow(rows, []any{&sr.UUID, &sr.Uploaded, &sr.Identifier, &sr.Title, &sr.ModelType}, func() error {
		res = append(res, sr)
		return nil
	})
	if err != nil {
		return err
	}

	if err = conn.Commit(ctx); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, res)
}

// @Summary		Run a model
// @Description	Runs a model with provided initial states and edge transfer functions.
// @Accept		json
// @Produce		json
// @Param		state		body		PostModelRunRequest	true	"state"
// @Success		200			{object}	PostModelRunResult
// @Failure		400			{object}	echo.HTTPError
// @Failure		404			{object}	echo.HTTPError
// @Failure		500			{object}	echo.HTTPError
// @Router		/model/run	[post]
func (self *State) PostModelRun(c echo.Context) error {
	ctx := c.Request().Context()

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	var mr PostModelRunRequest
	if err = json.Unmarshal(body, &mr); err != nil {
		c.Logger().Warn(err)
		return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: err.Error()})
	}

	r, err := self.Relations.Get(c.Logger())
	if err != nil {
		return err
	}

	if err := mr.Verify(&r); err != nil {
		c.Logger().Warn(err)
		return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: err.Error()})
	}

	query := `
		SELECT
			models.graph
		FROM
			models
		WHERE
			models.model_id = $1`

	conn, err := self.DBPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer conn.Rollback(ctx)

	var graph types.Graph
	if err = conn.QueryRow(ctx, query, mr.UUID).Scan(&graph); err != nil {
		if err == pgx.ErrNoRows {
			err = fmt.Errorf("model with uuid '%s' not found", mr.UUID)
			c.Logger().Warn(err)
			return c.JSON(http.StatusNotFound, echo.HTTPError{Message: err.Error()})
		} else {
			return err
		}
	}

	if err := conn.Commit(ctx); err != nil {
		return err
	}

	res, err := mr.Run(graph, &r)
	if err != nil {
		c.Logger().Warn(err)
		return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

// @Summary		Verify a model
// @Description	Verifies whether a model and its metadata are valid and complete.
// @Accept		json
// @Produce		json
// @Param		model			body		types.Model	true	"model"
// @Failure		400				{object}	echo.HTTPError
// @Failure		409				{object}	echo.HTTPError
// @Failure		500				{object}	echo.HTTPError
// @Router		/model/verify	[post]
func (self *State) PostModelVerify(c echo.Context) error {
	ctx := c.Request().Context()

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	var m types.Model
	if err = json.Unmarshal(body, &m); err != nil {
		c.Logger().Warn(err)
		return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: err.Error()})
	}

	if err = verifyModel(m, self, c.Logger(), ctx); err != nil {
		if _, ok := err.(types.Error); ok {
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

	if m.Parent != nil {
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

		var graph types.Graph
		var metadata []types.Metadata
		if err := conn.QueryRow(ctx, query, m.Parent).Scan(&graph, &metadata); err != nil {
			if err == pgx.ErrNoRows {
				err = fmt.Errorf("encountered non-existent parent model with uuid '%s'", *m.Parent)
				return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: err.Error()})
			} else {
				return err
			}
		}

		if slices.Equal(m.Graph.Nodes, graph.Nodes) && util.SlicesEqualUnordered(m.Graph.Relations, graph.Relations) && util.SlicesEqualUnordered(m.Metadata, metadata) {
			err = fmt.Errorf("model with uuid '%s' wasn't changed", *m.Parent)
			c.Logger().Warn(err)
			return c.JSON(http.StatusConflict, echo.HTTPError{Message: err.Error()})
		}
	}

	return c.JSON(http.StatusOK, struct{}{})
}
