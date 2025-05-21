-- DROP TABLE IF EXISTS models CASCADE;
CREATE TABLE IF NOT EXISTS models (
    model_id SERIAL PRIMARY KEY,
    model JSONB NOT NULL
);

-- DROP TABLE IF EXISTS metadata_values CASCADE;
CREATE TABLE IF NOT EXISTS metadata_values (
    value_id SERIAL PRIMARY KEY,
    model_id INT NOT NULL REFERENCES models (model_id),
    element_id INT NOT NULL,
    field_value TEXT NOT NULL,
    ontology TEXT
);
