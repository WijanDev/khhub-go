-- name: GetCongregation :one
SELECT id, name, number, midweek_day, weekend_day, updated_at
FROM congregation
WHERE id = 1;

-- name: UpdateCongregation :one
UPDATE congregation
SET name = $1,
    number = $2,
    midweek_day = $3,
    weekend_day = $4,
    updated_at = now()
WHERE id = 1
RETURNING id, name, number, midweek_day, weekend_day, updated_at;
