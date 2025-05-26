package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
)

func GetRelations(c echo.Context) error {
	cc := c.(*Context)
	relations, err := cc.Relations.Get(c.Logger())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, relations)
}

func GetMetadata(c echo.Context) error {
	cc := c.(*Context)
	metadata, err := cc.Metadata.Get(c.Logger())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, metadata)
}

func GetModel(c echo.Context) error {
	cc := c.(*Context)
	ctx := c.Request().Context()

	identifier := c.QueryParam("identifier")

	conn, err := cc.DbPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer conn.Rollback(ctx)

	query := `
		SELECT
			models.model,
			jsonb_agg(
				jsonb_build_object('kind', metadata.kind, 'value', metadata.value)
			)
		FROM
			models
			JOIN metadata ON models.model_id = metadata.model_id
		WHERE
			models.identifier = $1
		GROUP BY
			models.model`

	var model Model
	var metadata []Metadata
	err = conn.QueryRow(ctx, query, identifier).Scan(&model, &metadata)
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
	cc := c.(*Context)
	ctx := c.Request().Context()

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	var modelRequest ModelRequest
	err = json.Unmarshal(body, &modelRequest)
	if err != nil {
		c.Logger().Warn(err)
		return cc.String(http.StatusBadRequest, err.Error())
	}

	err = modelRequest.verify(cc)
	if err != nil {
		if e, ok := err.(GraphVerifyError); ok {
			c.Logger().Warn(e)
			return cc.String(http.StatusBadRequest, e.Error())
		} else {
			return err
		}
	}

	conn, err := cc.DbPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer conn.Rollback(ctx)

	const ERR_UNIQUE_VIOLATION = "23505"

	query := `
		INSERT INTO
			models (model, identifier)
		VALUES
			($1, $2)
		RETURNING
			model_id`

	var modelId int
	err = conn.QueryRow(ctx, query, modelRequest.Model, modelRequest.Identifier).Scan(&modelId)
	if err != nil {
		if e, ok := err.(*pgconn.PgError); ok && e.SQLState() == ERR_UNIQUE_VIOLATION {
			err = errors.New("encountered non-unique model identifier")
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
			($1, $2, $3)`

	var b pgx.Batch
	for _, metadata := range modelRequest.Metadata {
		b.Queue(query, modelId, metadata.Kind, metadata.Value)
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
	cc := c.(*Context)
	ctx := c.Request().Context()

	var query string
	var vs []any
	params := c.QueryParams()

	if len(params) == 0 {
		query = `
			SELECT
				models.identifier,
				metadata.value
			FROM
				models
				JOIN metadata ON models.model_id = metadata.model_id
			WHERE
				metadata.kind = 'title'`
	} else {
		e := map[string]struct{}{}
		i := 1
		q := []string{}
		for k, v := range params {
			e[k] = struct{}{}
			for _, ve := range v {
				vs = append(vs, k, "%"+ve+"%")
				q = append(q, fmt.Sprintf("(metadata.kind = $%d AND metadata.value ILIKE $%d)", i, i+1))
				i += 2
			}
		}
		vs = append(vs, len(e))

		query = fmt.Sprintf(`
			SELECT
				models.identifier,
				metadata.value
			FROM
				models
				JOIN metadata ON models.model_id = metadata.model_id
			WHERE
				metadata.kind = 'title'
				AND models.model_id IN (
					SELECT
						models.model_id
					FROM
						models
						JOIN metadata ON models.model_id = metadata.model_id
					WHERE
						%s
					GROUP BY
						models.model_id
					HAVING
						COUNT(DISTINCT metadata.kind) = $%d
				)`, strings.Join(q, " OR "), i)
	}

	conn, err := cc.DbPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer conn.Rollback(ctx)

	rows, err := conn.Query(ctx, query, vs...)
	if err != nil {
		return err
	}

	type SearchResult struct {
		Identifier string `json:"identifier"`
		Title      string `json:"title"`
	}

	var identifier, title string
	var idents []SearchResult
	_, err = pgx.ForEachRow(rows, []any{&identifier, &title}, func() error {
		idents = append(idents, SearchResult{identifier, title})
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
