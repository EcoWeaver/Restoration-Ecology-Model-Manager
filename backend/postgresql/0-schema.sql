CREATE TABLE IF NOT EXISTS models (
    model_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    uploaded TIMESTAMP NOT NULL DEFAULT now(),
    graph JSONB NOT NULL,
    parent UUID REFERENCES models (model_id)
);

CREATE INDEX IF NOT EXISTS models_parent_index ON models USING HASH(parent);

CREATE TABLE IF NOT EXISTS metadata (
    metadata_id SERIAL PRIMARY KEY,
    model_id UUID NOT NULL REFERENCES models (model_id),
    kind TEXT NOT NULL,
    value TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS metadata_model_id_index ON metadata USING HASH (model_id);
CREATE INDEX IF NOT EXISTS metadata_kind_index ON metadata USING HASH (kind);
