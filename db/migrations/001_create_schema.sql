CREATE TABLE IF NOT EXISTS mc_json_schema (
    schema_id serial PRIMARY KEY,
    schema json NOT NULL
);