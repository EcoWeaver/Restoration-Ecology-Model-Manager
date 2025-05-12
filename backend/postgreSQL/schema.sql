-- DROP TABLE IF EXISTS metadata_values CASCADE;

CREATE TABLE metadata_values (
  value_id SERIAL PRIMARY KEY,
  model_id INT,
  element_id INT,
  field_value TEXT,
  ontology TEXT
);
