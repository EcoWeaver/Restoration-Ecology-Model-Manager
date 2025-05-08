DROP TABLE IF EXISTS metadata_values;
DROP TABLE IF EXISTS element;
DROP TABLE IF EXISTS sub_category;
DROP TABLE IF EXISTS metadata_category;

CREATE TABLE metadata_category (
  metadata_category_id SERIAL PRIMARY KEY,
  main_category TEXT NOT NULL
);

CREATE TABLE sub_category (
  sub_category_id SERIAL PRIMARY KEY,
  metadata_category_id INT REFERENCES metadata_category(metadata_category_id) ON DELETE CASCADE,
  sub_category TEXT NOT NULL
);

CREATE TABLE element (
  element_id SERIAL PRIMARY KEY,
  sub_category_id INT REFERENCES sub_category(sub_category_id) ON DELETE CASCADE,
  element_name TEXT NOT NULL,
  element_description TEXT,
  is_mandatory BOOLEAN,
  allowed_ontology TEXT
);

CREATE TABLE metadata_values (
  value_id SERIAL PRIMARY KEY,
  model_id INT,
  element_id INT REFERENCES element(element_id) ON DELETE CASCADE,
  field_value TEXT,
  ontology TEXT
);
