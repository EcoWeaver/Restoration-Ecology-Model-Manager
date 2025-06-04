package main

import (
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

func (self *State) GetRelations(c echo.Context) error {
	relations, err := self.Relations.Get(c.Logger())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, relations)
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
	err := VerifyIdentifier(identifier)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	version := c.QueryParam("version")
	err = VerifyVersion(version)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	conn, err := self.DbPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer conn.Rollback(ctx)

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

	var graph Graph
	var metadata []Metadata
	err = conn.QueryRow(ctx, query, identifier, version).Scan(&graph, &metadata)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = fmt.Errorf("model with identifier '%s' at version '%s' not found", identifier, version)
			c.Logger().Warn(err)
			return c.String(http.StatusNotFound, err.Error())
		} else {
			return err
		}
	}

	err = conn.Commit(ctx)
	if err != nil {
		return err
	}

	r := ModelRequest{identifier, graph, version, metadata}
	err = r.verify(self, c.Logger())
	if err != nil {
		if _, ok := err.(ModelVerifyError); ok {
			err = fmt.Errorf("encountered error while verifying model with identifier '%s' at version '%s': %s", identifier, version, err.Error())
			c.Logger().Warn(err)
		} else {
			return err
		}
	}

	return c.JSON(http.StatusOK, r)
}

func (self *State) PutModel(c echo.Context) error {
	ctx := c.Request().Context()

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	var r ModelRequest
	err = json.Unmarshal(body, &r)
	if err != nil {
		c.Logger().Warn(err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	isCreate := r.Version == ""
	if isCreate {
		r.Version = "1"
	}

	err = r.verify(self, c.Logger())
	if err != nil {
		if _, ok := err.(ModelVerifyError); ok {
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

	if !isCreate {
		query := `
		UPDATE models
			SET num_derivatives = num_derivatives + 1
		WHERE
			identifier = $1 AND version = $2
		RETURNING
			num_derivatives`

		numDerivatives := 0
		err = conn.QueryRow(ctx, query, r.Identifier, r.Version).Scan(&numDerivatives)
		if err != nil {
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
	err = conn.QueryRow(ctx, query, r.Identifier, r.Graph, r.Version).Scan(&modelId)
	if err != nil {
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

	_, err = conn.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	err = conn.Commit(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, r)
}

func (self *State) SearchModel(c echo.Context) error {
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

	type SearchResult struct {
		Identifier string `json:"identifier"`
		Version    string `json:"version"`
		Title      string `json:"title"`
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

	err = conn.Commit(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, idents)
}
