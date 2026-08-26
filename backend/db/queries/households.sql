-- name: ListHouseholds :many
SELECT id, name, address, notes, created_at
FROM households
ORDER BY name;

-- name: GetHousehold :one
SELECT id, name, address, notes, created_at
FROM households
WHERE id = $1;

-- name: CreateHousehold :one
INSERT INTO households (name, address, notes)
VALUES ($1, $2, $3)
RETURNING id, name, address, notes, created_at;

-- name: UpdateHousehold :one
UPDATE households
SET name = $1, address = $2, notes = $3
WHERE id = $4
RETURNING id, name, address, notes, created_at;

-- name: DeleteHousehold :exec
DELETE FROM households WHERE id = $1;
