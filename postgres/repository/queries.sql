-- name: GetPerson :one
SELECT * FROM person WHERE id = $1;

-- name: AllPersons :many
SELECT * FROM person;

-- name: AddPerson :one
INSERT INTO person (name, age, location) VALUES ($1, $2, $3) RETURNING *;