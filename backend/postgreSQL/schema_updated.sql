DROP TABLE IF EXISTS metadata_values CASCADE;
DROP TABLE IF EXISTS element CASCADE; 
DROP TABLE IF EXISTS metadata_category CASCADE; 
DROP TABLE IF EXISTS sub_category CASCADE; 

CREATE TABLE metadata_values (
  value_id SERIAL PRIMARY KEY,
  model_id INT,
  element_id INT,
  field_value TEXT,
  ontology TEXT
);
