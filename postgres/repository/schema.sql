CREATE TABLE IF NOT EXISTS person
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    age INT NOT NULL,
    location GEOMETRY(POINT, 4326) NOT NULL
);

CREATE INDEX IF NOT EXISTS person_location_idx ON person USING GIST (location);
