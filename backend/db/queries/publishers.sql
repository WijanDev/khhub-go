-- name: CountPublishers :one
SELECT count(*) FROM publishers;

-- name: ListPublishers :many
SELECT
    p.id, p.household_id, p.first_name, p.last_name, p.gender, p.phone, p.email,
    p.baptism_date, p.started_preaching_date, p.spiritual_status,
    p.is_elder, p.is_ministerial_servant, p.is_regular_pioneer, p.is_special_pioneer,
    p.activity_status, p.created_at,
    h.name AS household_name
FROM publishers p
LEFT JOIN households h ON h.id = p.household_id
ORDER BY p.last_name, p.first_name;

-- name: GetPublisher :one
SELECT
    p.id, p.household_id, p.first_name, p.last_name, p.gender, p.phone, p.email,
    p.baptism_date, p.started_preaching_date, p.spiritual_status,
    p.is_elder, p.is_ministerial_servant, p.is_regular_pioneer, p.is_special_pioneer,
    p.activity_status, p.created_at,
    h.name AS household_name
FROM publishers p
LEFT JOIN households h ON h.id = p.household_id
WHERE p.id = $1;

-- name: CreatePublisher :one
INSERT INTO publishers (
    household_id, first_name, last_name, gender, phone, email,
    baptism_date, started_preaching_date, spiritual_status,
    is_elder, is_ministerial_servant, is_regular_pioneer, is_special_pioneer
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
)
RETURNING id, household_id, first_name, last_name, gender, phone, email,
    baptism_date, started_preaching_date, spiritual_status,
    is_elder, is_ministerial_servant, is_regular_pioneer, is_special_pioneer,
    activity_status, created_at;

-- name: UpdatePublisher :one
UPDATE publishers SET
    household_id = $2,
    first_name = $3,
    last_name = $4,
    gender = $5,
    phone = $6,
    email = $7,
    baptism_date = $8,
    started_preaching_date = $9,
    spiritual_status = $10,
    is_elder = $11,
    is_ministerial_servant = $12,
    is_regular_pioneer = $13,
    is_special_pioneer = $14
WHERE id = $1
RETURNING id, household_id, first_name, last_name, gender, phone, email,
    baptism_date, started_preaching_date, spiritual_status,
    is_elder, is_ministerial_servant, is_regular_pioneer, is_special_pioneer,
    activity_status, created_at;

-- name: DeletePublisher :exec
DELETE FROM publishers WHERE id = $1;

-- name: UpdatePublisherActivity :exec
UPDATE publishers SET activity_status = $2 WHERE id = $1;

-- name: ListReportingPublishers :many
SELECT
    p.id, p.household_id, p.first_name, p.last_name, p.gender, p.phone, p.email,
    p.baptism_date, p.started_preaching_date, p.spiritual_status,
    p.is_elder, p.is_ministerial_servant, p.is_regular_pioneer, p.is_special_pioneer,
    p.activity_status, p.created_at
FROM publishers p
WHERE p.spiritual_status IN ('unbaptized_publisher', 'publisher')
ORDER BY p.last_name, p.first_name;
