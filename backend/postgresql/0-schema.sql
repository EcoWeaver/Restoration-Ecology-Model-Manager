CREATE TABLE IF NOT EXISTS models (
    model_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created TIMESTAMP NOT NULL DEFAULT now(),
    graph JSONB NOT NULL,
    parent UUID REFERENCES models (model_id)
);

CREATE TABLE IF NOT EXISTS metadata (
    metadata_id SERIAL PRIMARY KEY,
    model_id UUID NOT NULL REFERENCES models (model_id),
    kind TEXT NOT NULL,
    value TEXT NOT NULL
);
