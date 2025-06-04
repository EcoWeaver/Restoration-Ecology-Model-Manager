CREATE TABLE IF NOT EXISTS models (
    model_id SERIAL PRIMARY KEY,
    identifier TEXT NOT NULL,
    graph JSONB NOT NULL,
    version TEXT NOT NULL,
    num_derivatives INT NOT NULL DEFAULT 0,
    UNIQUE (identifier, version)
);

CREATE TABLE IF NOT EXISTS metadata (
    metadata_id SERIAL PRIMARY KEY,
    model_id INT NOT NULL REFERENCES models (model_id),
    kind TEXT NOT NULL,
    value TEXT NOT NULL
);
