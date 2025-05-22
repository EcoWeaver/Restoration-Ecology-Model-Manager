CREATE TABLE IF NOT EXISTS models (
    model_id SERIAL PRIMARY KEY,
    model JSONB NOT NULL,
    identifier TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS metadata (
    metadata_id SERIAL PRIMARY KEY,
    model_id INT NOT NULL REFERENCES models (model_id),
    kind TEXT NOT NULL,
    value TEXT NOT NULL
);
