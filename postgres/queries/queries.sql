-- name: GetPerson :one
SELECT * FROM person WHERE id = $1;

-- name: AllPersons :many
SELECT * FROM person;

-- name: AddPerson :one
INSERT INTO person (name, age, location) VALUES ($1, $2, $3) RETURNING *;

-- name: UpdatePerson :one
UPDATE person SET name = $1, age = $2, location = $3 WHERE id = $4 RETURNING *;

-- name: UpsertPerson :one
INSERT INTO person (name, age, location) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET name = $1, age = $2, location = $3 RETURNING *;

-- name: DeletePerson :one
DELETE FROM person WHERE id = $1 RETURNING *;