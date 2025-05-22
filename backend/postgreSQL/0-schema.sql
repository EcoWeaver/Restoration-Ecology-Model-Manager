CREATE TABLE IF NOT EXISTS models (
    model_id SERIAL PRIMARY KEY,
    identifier TEXT UNIQUE NOT NULL,
    related_identifier TEXT, 
    model JSONB NOT NULL
);

CREATE TABLE IF NOT EXISTS metadata_values (
    value_id SERIAL PRIMARY KEY,
    model_id INT NOT NULL REFERENCES models (model_id),
    element_id INT NOT NULL,
    field_value TEXT NOT NULL,
    ontology TEXT
);

CREATE TABLE IF NOT EXISTS controlled_terms (
    controlled_terms_id SERIAL PRIMARY KEY, 
    element_id INT NOT NULL, 
    value TEXT NOT NULL, 
    ontology_URI TEXT 
); 