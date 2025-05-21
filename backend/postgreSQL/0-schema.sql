CREATE TABLE IF NOT EXISTS models (
    model_id SERIAL PRIMARY KEY,
    identifier TEXT UNIQUE NOT NULL,
    model JSONB NOT NULL
);

CREATE TABLE IF NOT EXISTS metadata_values (
    value_id SERIAL PRIMARY KEY,
    model_id INT NOT NULL REFERENCES models (model_id),
    element_id INT NOT NULL,
    field_value TEXT NOT NULL,
    ontology TEXT
);
