package main

import (
	"backend/types"
	"backend/verify"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"
)

type State struct {
	DBPool       *pgxpool.Pool
	Metadata     HotReload[MetadataConfig]
	Relations    HotReload[RelationsConfig]
	DOICache     verify.Cache
	ORCIDCache   verify.Cache
	SpeciesCache verify.Cache
	OLS4Caches   sync.Map
}

type GetModelResult struct {
	UUID     string      `json:"uuid"`
	Uploaded time.Time   `json:"uploaded"`
	Model    types.Model `json:"model"`
	Parent   *Parent     `json:"parent"`
}

type PutModelRequest struct {
	Model  types.Model `json:"model"`
	Parent *string     `json:"parent"`
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

type PostSearchRequest struct {
	Metadata   []types.Metadata `json:"modelMetadata"`
	NodeBodies []string         `json:"nodeBodies"`
}

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
func (self *State) GetMetadata(c *echo.Context) error {
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
func (self *State) GetModel(c *echo.Context) error {
	ctx := c.Request().Context()

	modelUUID := c.QueryParam("uuid")
	if err := verify.UUID(modelUUID); err != nil {
		return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: err.Error()})
	}

	query := `
		SELECT
			models.uploaded,
			models.parent
		FROM
			models
		WHERE
			models.model_id = $1`

	conn, err := self.DBPool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	var uploaded time.Time
	var parentUUID *string
	if err = conn.QueryRow(ctx, query, modelUUID).Scan(&uploaded, &parentUUID); err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.HTTPError{Message: fmt.Sprintf("model with uuid '%s' not found", modelUUID)})
		} else {
			return err
		}
	}

	model, err := getModelFromDB(ctx, conn.Conn(), modelUUID)
	if err != nil {
		return err
	}

	var parent *Parent
	if parentUUID != nil {
		parent = &Parent{UUID: *parentUUID}

		query = `
		SELECT
			model_metadata.value
		FROM
			model_metadata
		WHERE
			model_metadata.model_id = $1 AND model_metadata.kind = 'title'`

		if err := conn.QueryRow(ctx, query, *parentUUID).Scan(&parent.Title); err != nil {
			return err
		}
	}

	go func() {
		if err, ok := verifyModel(model, self, c.Logger(), context.Background()).(types.Error); ok {
			c.Logger().Warn(fmt.Sprintf("encountered error while verifying model with uuid '%s': %s", modelUUID, err.Error()))
		}
	}()

	return c.JSON(http.StatusOK, GetModelResult{modelUUID, uploaded, model, parent})
}

// @Summary		Save a model
// @Description	Saves a model. Returns the UUID assigned to that model.
// @Accept		json
// @Produce		json
// @Param		model	body		PutModelRequest	true	"model"
// @Success		201		{object}	PutModelResult
// @Failure		400		{object}	echo.HTTPError
// @Failure		409		{object}	echo.HTTPError
// @Failure		500		{object}	echo.HTTPError
// @Router		/model	[put]
func (self *State) PutModel(c *echo.Context) error {
	ctx := c.Request().Context()

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	var m PutModelRequest
	if err = json.Unmarshal(body, &m); err != nil {
		return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: err.Error()})
	}

	if err = m.verify(self, c.Logger(), ctx); err != nil {
		if _, ok := err.(types.Error); ok {
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
		if err := checkModelExistsInDB(ctx, conn.Conn(), *m.Parent); err != nil {
			if err == pgx.ErrNoRows {
				return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: fmt.Sprintf("encountered non-existent parent model with uuid '%s'", *m.Parent)})
			} else {
				return err
			}
		}

		model, err := getModelFromDB(ctx, conn.Conn(), *m.Parent)
		if err != nil {
			return err
		}

		if modelsEqual(m.Model, model) {
			return c.JSON(http.StatusConflict, echo.HTTPError{Message: fmt.Sprintf("model with uuid '%s' wasn't changed", *m.Parent)})
		}
	}

	query := `
		INSERT INTO
			models (parent)
		VALUES
			($1)
		RETURNING
			model_id`

	var res PutModelResult
	if err = conn.QueryRow(ctx, query, m.Parent).Scan(&res.UUID); err != nil {
		return err
	}

	if err = insertModelIntoDB(ctx, conn, res.UUID, m.Model); err != nil {
		return err
	}

	if err = conn.Commit(ctx); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, res)
}

// @Summary		Get the relations
// @Description	Gets the list of currently accepted relation kinds and types.
// @Produce		json
// @Success		200			{object}	RelationsConfig
// @Failure		500			{object}	echo.HTTPError
// @Router		/relations	[get]
func (self *State) GetRelations(c *echo.Context) error {
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
// @Param		queries	body		PostSearchRequest	true	"queries"
// @Success		200		{object}	[]PostSearchResult
// @Failure		400		{object}	echo.HTTPError
// @Failure		500		{object}	echo.HTTPError
// @Router		/search	[post]
func (self *State) PostSearch(c *echo.Context) error {
	ctx := c.Request().Context()

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	var r PostSearchRequest
	if err = json.Unmarshal(body, &r); err != nil {
		return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: err.Error()})
	}

	m, err := self.Metadata.Get(c.Logger())
	if err != nil {
		return err
	}

	qs := map[string][]string{}
	for _, mq := range r.Metadata {
		if _, ok := qs[mq.Kind]; !ok && !m.IsMetadataKind(mq.Kind) {
			return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: fmt.Sprintf("encountered invalid metadata kind '%s'", mq.Kind)})
		}
		qs[mq.Kind] = append(qs[mq.Kind], mq.Value)
	}

	query := `
		SELECT
			models.model_id,
			models.uploaded,
			MAX(CASE WHEN model_metadata.kind = 'identifier' THEN model_metadata.value END),
			MAX(CASE WHEN model_metadata.kind = 'title' THEN model_metadata.value END),
			MAX(CASE WHEN model_metadata.kind = 'model type' THEN model_metadata.value END)
		FROM
			models
			JOIN model_metadata ON models.model_id = model_metadata.model_id`

	var args []any
	if len(qs) != 0 {
		var conds []string
		for k, v := range qs {
			args = append(args, k)
			conj := fmt.Sprintf("(model_metadata.kind = $%d AND (", len(args))

			var disj []string
			for _, ve := range v {
				args = append(args, "%"+ve+"%")
				disj = append(disj, fmt.Sprintf("model_metadata.value ILIKE $%d", len(args)))
			}

			conds = append(conds, conj+strings.Join(disj, " OR ")+"))")
		}

		innerQuery := `
				SELECT
					models.model_id
				FROM
					models
					JOIN model_metadata ON models.model_id = model_metadata.model_id
				WHERE
					` + strings.Join(conds, `
					OR `)

		if len(qs) != 1 {
			args = append(args, len(qs))

			innerQuery += fmt.Sprintf(`
				GROUP BY
					models.model_id
				HAVING
					COUNT(DISTINCT model_metadata.kind) = $%d`, len(args))
		}

		query += `
		WHERE
			models.model_id IN (` + innerQuery + `
			)`
	}

	if len(r.NodeBodies) != 0 {
		var disj []string
		for _, ve := range r.NodeBodies {
			args = append(args, "%"+ve+"%")
			disj = append(disj, fmt.Sprintf("nodes.body ILIKE $%d", len(args)))
		}

		innerQuery := `
				SELECT
					models.model_id
				FROM
					models
					JOIN nodes ON models.model_id = nodes.model_id
				WHERE
					` + strings.Join(disj, " OR ")

		if len(qs) != 0 {
			query += `
		AND`
		} else {
			query += `
		WHERE`
		}

		query += `
			models.model_id IN (` + innerQuery + `
			)`
	}

	query += `
		GROUP BY
			models.model_id,
			models.uploaded`

	conn, err := self.DBPool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	var sr PostSearchResult
	var res []PostSearchResult
	_, err = pgx.ForEachRow(rows, []any{&sr.UUID, &sr.Uploaded, &sr.Identifier, &sr.Title, &sr.ModelType}, func() error {
		res = append(res, sr)
		return nil
	})
	if err != nil {
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
func (self *State) PostModelRun(c *echo.Context) error {
	ctx := c.Request().Context()

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	var mr PostModelRunRequest
	if err = json.Unmarshal(body, &mr); err != nil {
		return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: err.Error()})
	}

	if err := mr.Verify(); err != nil {
		return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: err.Error()})
	}

	conn, err := self.DBPool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	var one int
	if err := conn.QueryRow(ctx, "SELECT 1 FROM models WHERE models.model_id = $1", mr.UUID).Scan(&one); err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.HTTPError{Message: fmt.Sprintf("model with uuid '%s' not found", mr.UUID)})
		} else {
			return err
		}
	}

	graph, err := getModelGraphFromDB(ctx, conn.Conn(), mr.UUID)
	if err != nil {
		return err
	}

	r, err := self.Relations.Get(c.Logger())
	if err != nil {
		return err
	}

	s, err := mr.Build(graph, r)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: err.Error()})
	}

	res, err := s.Run()
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

// @Summary		Verify a model
// @Description	Verifies whether a model and its metadata are valid and complete.
// @Accept		json
// @Produce		json
// @Param		model			body		PutModelRequest	true	"model"
// @Failure		400				{object}	echo.HTTPError
// @Failure		409				{object}	echo.HTTPError
// @Failure		500				{object}	echo.HTTPError
// @Router		/model/verify	[post]
func (self *State) PostModelVerify(c *echo.Context) error {
	ctx := c.Request().Context()

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	var m PutModelRequest
	if err = json.Unmarshal(body, &m); err != nil {
		return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: err.Error()})
	}

	if err = m.verify(self, c.Logger(), ctx); err != nil {
		if _, ok := err.(types.Error); ok {
			return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: err.Error()})
		} else {
			return err
		}
	}

	if m.Parent != nil {
		conn, err := self.DBPool.Acquire(ctx)
		if err != nil {
			return err
		}
		defer conn.Release()

		if err := checkModelExistsInDB(ctx, conn.Conn(), *m.Parent); err != nil {
			if err == pgx.ErrNoRows {
				return c.JSON(http.StatusBadRequest, echo.HTTPError{Message: fmt.Sprintf("encountered non-existent parent model with uuid '%s'", *m.Parent)})
			} else {
				return err
			}
		}

		model, err := getModelFromDB(ctx, conn.Conn(), *m.Parent)
		if err != nil {
			return err
		}

		if modelsEqual(m.Model, model) {
			return c.JSON(http.StatusConflict, echo.HTTPError{Message: fmt.Sprintf("model with uuid '%s' wasn't changed", *m.Parent)})
		}
	}

	return c.JSON(http.StatusOK, struct{}{})
}

func checkModelExistsInDB(ctx context.Context, conn *pgx.Conn, modelUUID string) error {
	var one int
	if err := conn.QueryRow(ctx, "SELECT 1 FROM models WHERE models.model_id = $1", modelUUID).Scan(&one); err != nil {
		return err
	}

	return nil
}

func getModelFromDB(ctx context.Context, conn *pgx.Conn, modelUUID string) (types.Model, error) {
	graph, err := getModelGraphFromDB(ctx, conn, modelUUID)
	if err != nil {
		return types.Model{}, err
	}

	metadata, err := getModelMetadataFromDB(ctx, conn, modelUUID)
	if err != nil {
		return types.Model{}, err
	}

	return types.Model{Graph: graph, Metadata: metadata}, nil
}

func getModelGraphFromDB(ctx context.Context, conn *pgx.Conn, modelUUID string) (types.Graph, error) {
	var graph types.Graph

	query := `
		SELECT
			nodes.node_id,
			nodes.body,
			nodes.color,
			nodes.x,
			nodes.y
		FROM
			nodes
		WHERE
			nodes.model_id = $1
		ORDER BY
			nodes.ord`

	rows, err := conn.Query(ctx, query, modelUUID)
	if err != nil {
		return types.Graph{}, err
	}
	defer rows.Close()

	var id int
	var n types.Node
	ids := map[int]int{}
	_, err = pgx.ForEachRow(rows, []any{&id, &n.Body, &n.Color, &n.Position.X, &n.Position.Y}, func() error {
		ids[id] = len(graph.Nodes)
		graph.Nodes = append(graph.Nodes, n)
		return nil
	})
	if err != nil {
		return types.Graph{}, err
	}

	query = `
		SELECT
			relations.from_id,
			relations.to_id,
			relations.kind,
			relations.r,
			relations.phi
		FROM
			relations
		WHERE
			relations.model_id = $1`

	relRows, err := conn.Query(ctx, query, modelUUID)
	if err != nil {
		return types.Graph{}, err
	}
	defer relRows.Close()

	var r types.Relation
	_, err = pgx.ForEachRow(relRows, []any{&r.From, &r.To, &r.Kind, &r.Position.R, &r.Position.Phi}, func() error {
		r.From = ids[r.From]
		r.To = ids[r.To]
		graph.Relations = append(graph.Relations, r)
		return nil
	})
	if err != nil {
		return types.Graph{}, err
	}

	return graph, nil
}

func getModelMetadataFromDB(ctx context.Context, conn *pgx.Conn, modelUUID string) ([]types.Metadata, error) {
	query := `
		SELECT
			model_metadata.kind,
			model_metadata.value
		FROM
			model_metadata
		WHERE
			model_metadata.model_id = $1`

	rows, err := conn.Query(ctx, query, modelUUID)
	if err != nil {
		return []types.Metadata{}, err
	}
	defer rows.Close()

	var me types.Metadata
	var metadata []types.Metadata
	_, err = pgx.ForEachRow(rows, []any{&me.Kind, &me.Value}, func() error {
		metadata = append(metadata, me)
		return nil
	})
	if err != nil {
		return []types.Metadata{}, err
	}

	return metadata, nil
}

func insertModelIntoDB(ctx context.Context, conn pgx.Tx, modelUUID string, g types.Model) error {
	if err := insertModelGraphIntoDB(ctx, conn, modelUUID, g.Graph); err != nil {
		return err
	}

	if err := insertModelMetadataIntoDB(ctx, conn, modelUUID, g.Metadata); err != nil {
		return err
	}

	return nil
}

func insertModelGraphIntoDB(ctx context.Context, conn pgx.Tx, modelUUID string, g types.Graph) error {
	if len(g.Nodes) == 0 {
		return nil
	}

	args := []any{modelUUID}
	var vals []string
	for i, n := range g.Nodes {
		args = append(args, i, n.Body, n.Color, n.Position.X, n.Position.Y)
		vals = append(vals,
			fmt.Sprintf("($1, $%d, $%d, $%d, $%d, $%d)",
				len(args)-4, len(args)-3, len(args)-2, len(args)-1, len(args)),
		)
	}

	query := `
		WITH inserted AS (
			INSERT INTO
				nodes (model_id, ord, body, color, x, y)
			VALUES
				` + strings.Join(vals, `,
				`) + `
			RETURNING
				nodes.node_id,
				nodes.ord
		)
		SELECT
			inserted.node_id
		FROM
			inserted
		ORDER BY
			inserted.ord`

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	var id int
	args = []any{modelUUID}
	_, err = pgx.ForEachRow(rows, []any{&id}, func() error {
		args = append(args, id)
		return nil
	})
	if err != nil {
		return err
	}

	if len(g.Relations) == 0 {
		return nil
	}

	vals = nil
	for _, r := range g.Relations {
		args = append(args, r.Kind, r.Position.R, r.Position.Phi)
		vals = append(vals, fmt.Sprintf("($1, $%d, $%d, $%d, $%d, $%d)", r.From+2, r.To+2, len(args)-2, len(args)-1, len(args)))
	}

	query = `
		INSERT INTO
			relations (model_id, from_id, to_id, kind, r, phi)
		VALUES
			` + strings.Join(vals, `,
			`)

	if _, err = conn.Exec(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func insertModelMetadataIntoDB(ctx context.Context, conn pgx.Tx, modelUUID string, metadata []types.Metadata) error {
	args := []any{modelUUID}
	var vals []string
	for _, me := range metadata {
		args = append(args, me.Kind, me.Value)
		vals = append(vals, fmt.Sprintf("($1, $%d, $%d)", len(args)-1, len(args)))
	}

	query := `
		INSERT INTO
			model_metadata (model_id, kind, value)
		VALUES
			` + strings.Join(vals, `,
			`)

	if _, err := conn.Exec(ctx, query, args...); err != nil {
		return err
	}

	return nil
}
