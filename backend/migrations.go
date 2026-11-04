package main

import (
	"backend/types"
	"context"
	"log/slog"
	"reflect"
	"runtime"

	"github.com/jackc/pgx/v5"
)

func (self *State) runMigrations(ctx context.Context, l *slog.Logger) error {
	migrations := []func(context.Context, pgx.Tx) error{
		schema, migrateModelsGraphJSON, migrateModelsParent, migrateRelationMetadataAddOrd,
	}

	for _, migration := range migrations {
		if err := self.runMigration(ctx, migration, l); err != nil {
			return err
		}
	}

	return nil
}

func (self *State) runMigration(ctx context.Context, f func(context.Context, pgx.Tx) error, l *slog.Logger) error {
	id := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()

	conn, err := self.DBPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer conn.Rollback(ctx)

	if _, err := conn.Exec(ctx, "CREATE TABLE IF NOT EXISTS migrations (id TEXT PRIMARY KEY)"); err != nil {
		return err
	}

	var dummy int
	if err := conn.QueryRow(ctx, "SELECT 1 FROM migrations WHERE migrations.id = $1", id).Scan(&dummy); err != pgx.ErrNoRows {
		return err
	}

	l.Info("running migration <" + id + ">")

	if err := f(ctx, conn); err != nil {
		return err
	}

	if _, err := conn.Exec(ctx, "INSERT INTO migrations (id) VALUES ($1)", id); err != nil {
		return err
	}

	if err := conn.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func schema(ctx context.Context, conn pgx.Tx) error {
	query := `
		CREATE TABLE IF NOT EXISTS models (
			model_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			uploaded TIMESTAMP NOT NULL DEFAULT now(),
			graph JSONB NOT NULL,
			parent UUID REFERENCES models (model_id)
		);

		CREATE TABLE IF NOT EXISTS metadata (
			metadata_id SERIAL PRIMARY KEY,
			model_id UUID NOT NULL REFERENCES models (model_id),
			kind TEXT NOT NULL,
			value TEXT NOT NULL
		);`

	if _, err := conn.Exec(ctx, query); err != nil {
		return err
	}

	return nil
}

func migrateModelsGraphJSON(ctx context.Context, conn pgx.Tx) error {
	query := `
		ALTER TABLE metadata RENAME TO model_metadata;

		CREATE TABLE IF NOT EXISTS nodes (
			node_id SERIAL PRIMARY KEY,
			model_id UUID NOT NULL REFERENCES models (model_id),
			ord INT NOT NULL,
			body TEXT NOT NULL,
			color TEXT NOT NULL,
			x FLOAT NOT NULL,
			y FLOAT NOT NULL,
			UNIQUE (model_id, ord)
		);

		CREATE TABLE IF NOT EXISTS node_metadata (
			metadata_id SERIAL PRIMARY KEY,
			model_id UUID NOT NULL REFERENCES models (model_id),
			node_id INT NOT NULL REFERENCES nodes (node_id),
			kind TEXT NOT NULL,
			value TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS relations (
			relation_id SERIAL PRIMARY KEY,
			model_id UUID NOT NULL REFERENCES models (model_id),
			from_id INT NOT NULL REFERENCES nodes (node_id),
			to_id INT NOT NULL REFERENCES nodes (node_id),
			ord INT NOT NULL,
			kind TEXT NOT NULL,
			r FLOAT NOT NULL,
			phi FLOAT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS relation_metadata (
			metadata_id SERIAL PRIMARY KEY,
			model_id UUID NOT NULL REFERENCES models (model_id),
			relation_id INT NOT NULL REFERENCES relations (relation_id),
			kind TEXT NOT NULL,
			value TEXT NOT NULL
		);

		CREATE INDEX IF NOT EXISTS model_metadata_model_index ON model_metadata (model_id);
		CREATE INDEX IF NOT EXISTS model_metadata_kind_index ON model_metadata (kind);

		CREATE INDEX IF NOT EXISTS nodes_model_index ON nodes (model_id);

		CREATE INDEX IF NOT EXISTS node_metadata_model_index ON node_metadata (model_id);
		CREATE INDEX IF NOT EXISTS node_metadata_kind_index ON node_metadata (kind);

		CREATE INDEX IF NOT EXISTS relations_model_index ON relations (model_id);

		CREATE INDEX IF NOT EXISTS relation_metadata_model_index ON relation_metadata (model_id);
		CREATE INDEX IF NOT EXISTS relation_metadata_kind_index ON relation_metadata (kind);`

	if _, err := conn.Exec(ctx, query); err != nil {
		return err
	}

	rows, err := conn.Query(ctx, "SELECT models.model_id, models.graph FROM models")
	if err != nil {
		return err
	}
	defer rows.Close()

	type Model struct {
		modelUUID string
		graph     types.Graph
	}

	var model Model
	var models []Model
	_, err = pgx.ForEachRow(rows, []any{&model.modelUUID, &model.graph}, func() error {
		models = append(models, model)
		return nil
	})
	if err != nil {
		return err
	}

	for _, model := range models {
		if err = insertModelIntoDB(ctx, conn, model.modelUUID, types.Model{Graph: model.graph, Metadata: []types.Metadata{}}); err != nil {
			return err
		}
	}

	if _, err := conn.Exec(ctx, "ALTER TABLE models DROP COLUMN graph"); err != nil {
		return err
	}

	return nil
}

func migrateModelsParent(ctx context.Context, conn pgx.Tx) error {
	query := `
		INSERT INTO
			model_metadata (model_id, kind, value)
		SELECT
			models.model_id,
			'derived from',
			models.parent
		FROM
			models
		WHERE
			models.parent IS NOT NULL;

		ALTER TABLE models DROP COLUMN parent;`

	if _, err := conn.Exec(ctx, query); err != nil {
		return err
	}

	return nil
}

func migrateRelationMetadataAddOrd(ctx context.Context, conn pgx.Tx) error {
	if _, err := conn.Exec(ctx, "ALTER TABLE relations ADD COLUMN IF NOT EXISTS ord INT NOT NULL"); err != nil {
		return err
	}

	return nil
}
