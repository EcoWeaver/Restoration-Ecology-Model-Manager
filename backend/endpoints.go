package main

import (
	"backend/verify"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type State struct {
	DbPool    *pgxpool.Pool
	Metadata  LiveMetadataConfig
	Relations LiveRelationsConfig
}

type SearchResult struct {
	Identifier string `json:"identifier"`
	Version    string `json:"version"`
	Title      string `json:"title"`
}

func (self *State) GetMetadata(c echo.Context) error {
	metadata, err := self.Metadata.Get(c.Logger())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, metadata)
}

func (self *State) GetModel(c echo.Context) error {
	ctx := c.Request().Context()

	identifier := c.QueryParam("identifier")
	version := c.QueryParam("version")

	query := `
		SELECT
			models.graph,
			jsonb_agg(
				jsonb_build_object('kind', metadata.kind, 'value', metadata.value)
			)
		FROM
			models
			JOIN metadata ON models.model_id = metadata.model_id
		WHERE
			models.identifier = $1 AND models.version = $2
		GROUP BY
			models.graph`

	conn, err := self.DbPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer conn.Rollback(ctx)

	var graph Graph
	var metadata []Metadata
	if err = conn.QueryRow(ctx, query, identifier, version).Scan(&graph, &metadata); err != nil {
		if err == pgx.ErrNoRows {
			err = fmt.Errorf("model with identifier '%s' at version '%s' not found", identifier, version)
			c.Logger().Warn(err)
			return c.String(http.StatusNotFound, err.Error())
		} else {
			return err
		}
	}

	if err = conn.Commit(ctx); err != nil {
		return err
	}

	r := ModelRequest{identifier, graph, version, metadata}

	go func() {
		if err := r.verify(self, c.Logger()); err != nil {
			c.Logger().Warn(fmt.Sprintf("encountered error while verifying model with identifier '%s' at version '%s': %s", identifier, version, err))
		}
	}()

	return c.JSON(http.StatusOK, r)
}

func (self *State) PutModel(c echo.Context) error {
	ctx := c.Request().Context()

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	var r ModelRequest
	if err = json.Unmarshal(body, &r); err != nil {
		c.Logger().Warn(err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	isDerivative := r.Version != ""
	if !isDerivative {
		r.Version = "1"
	}

	if err = r.verify(self, c.Logger()); err != nil {
		if _, ok := err.(verify.VerifyError); ok {
			c.Logger().Warn(err)
			return c.String(http.StatusBadRequest, err.Error())
		} else {
			return err
		}
	}

	conn, err := self.DbPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer conn.Rollback(ctx)

	if isDerivative {
		query := `
		UPDATE models
			SET num_derivatives = num_derivatives + 1
		WHERE
			identifier = $1 AND version = $2
		RETURNING
			num_derivatives`

		var numDerivatives int
		if err = conn.QueryRow(ctx, query, r.Identifier, r.Version).Scan(&numDerivatives); err != nil {
			if err == pgx.ErrNoRows {
				err = fmt.Errorf("model with identifier '%s' at version '%s' not found", r.Identifier, r.Version)
				c.Logger().Warn(err)
				return c.String(http.StatusNotFound, err.Error())
			} else {
				return err
			}
		}

		r.Version += fmt.Sprintf(".%d", numDerivatives)
	}

	const ERR_UNIQUE_VIOLATION = "23505"

	query := `
		INSERT INTO
			models (identifier, graph, version)
		VALUES
			($1, $2, $3)
		RETURNING
			model_id`

	var modelId int
	if err = conn.QueryRow(ctx, query, r.Identifier, r.Graph, r.Version).Scan(&modelId); err != nil {
		if e, ok := err.(*pgconn.PgError); ok && e.SQLState() == ERR_UNIQUE_VIOLATION {
			err = fmt.Errorf("encountered non-unique model identifier '%s' at version '%s'", r.Identifier, r.Version)
			c.Logger().Warn(err)
			return c.String(http.StatusBadRequest, err.Error())
		} else {
			return err
		}
	}

	query = `
		INSERT INTO
			metadata (model_id, kind, value)
		VALUES
			`

	args := []any{modelId}
	ins := []string{}
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

func (self *State) GetRelations(c echo.Context) error {
	relations, err := self.Relations.Get(c.Logger())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, relations)
}

func (self *State) GetSearch(c echo.Context) error {
	ctx := c.Request().Context()

	metadata, err := self.Metadata.Get(c.Logger())
	if err != nil {
		return err
	}

	query := `
		SELECT
			models.identifier,
			models.version,
			metadata.value
		FROM
			models
			JOIN metadata ON models.model_id = metadata.model_id
		WHERE
			metadata.kind = 'title'`
	var args []any

	params := c.QueryParams()
	if len(params) != 0 {
		conds := []string{}
		for k, v := range params {
			if !metadata.IsMetadataKind(k) {
				err = fmt.Errorf("encountered invalid metadata kind '%s'", k)
				c.Logger().Warn(err)
				return c.String(http.StatusBadRequest, err.Error())
			}

			args = append(args, k)
			conj := fmt.Sprintf("(metadata.kind = $%d AND (", len(args))

			disj := []string{}
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
			AND models.model_id IN (` + innerQuery + `
			)`
	}

	conn, err := self.DbPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer conn.Rollback(ctx)

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return err
	}

	var identifier, version, title string
	var idents []SearchResult
	_, err = pgx.ForEachRow(rows, []any{&identifier, &version, &title}, func() error {
		idents = append(idents, SearchResult{identifier, version, title})
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

func (self *State) GetVersions(c echo.Context) error {
	ctx := c.Request().Context()

	identifier := c.QueryParam("identifier")

	query := `
		SELECT
			models.version,
			metadata.value
		FROM
			models
			JOIN metadata ON models.model_id = metadata.model_id
		WHERE
			models.identifier = $1 AND metadata.kind = 'title'`

	conn, err := self.DbPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer conn.Rollback(ctx)

	rows, err := conn.Query(ctx, query, identifier)
	if err != nil {
		return err
	}

	var version, title string
	var idents []SearchResult
	_, err = pgx.ForEachRow(rows, []any{&version, &title}, func() error {
		idents = append(idents, SearchResult{identifier, version, title})
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
