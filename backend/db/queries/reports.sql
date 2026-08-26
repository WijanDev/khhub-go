-- name: ListReportsForMonth :many
SELECT
    r.id, p.id AS publisher_id, r.year, r.month, r.shared_in_ministry, r.bible_studies,
    r.hours, r.auxiliary_pioneer, r.late, r.remarks, r.created_at, r.updated_at,
    p.first_name, p.last_name, p.spiritual_status,
    p.is_regular_pioneer, p.is_special_pioneer
FROM publishers p
LEFT JOIN field_service_reports r
    ON r.publisher_id = p.id AND r.year = sqlc.arg(year) AND r.month = sqlc.arg(month)
WHERE p.spiritual_status IN ('unbaptized_publisher', 'publisher')
ORDER BY p.last_name, p.first_name;

-- name: UpsertReport :one
INSERT INTO field_service_reports (
    publisher_id, year, month, shared_in_ministry, bible_studies,
    hours, auxiliary_pioneer, late, remarks, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, now()
)
ON CONFLICT (publisher_id, year, month) DO UPDATE SET
    shared_in_ministry = EXCLUDED.shared_in_ministry,
    bible_studies = EXCLUDED.bible_studies,
    hours = EXCLUDED.hours,
    auxiliary_pioneer = EXCLUDED.auxiliary_pioneer,
    late = EXCLUDED.late,
    remarks = EXCLUDED.remarks,
    updated_at = now()
RETURNING id, publisher_id, year, month, shared_in_ministry, bible_studies,
    hours, auxiliary_pioneer, late, remarks, created_at, updated_at;

-- name: ListSharesForPublisher :many
SELECT year, month, shared_in_ministry
FROM field_service_reports
WHERE publisher_id = $1
  AND (year > sqlc.arg(from_year) OR (year = sqlc.arg(from_year) AND month >= sqlc.arg(from_month)))
  AND (year < sqlc.arg(to_year) OR (year = sqlc.arg(to_year) AND month <= sqlc.arg(to_month)))
ORDER BY year, month;

-- name: GetPublisherForReport :one
SELECT id, started_preaching_date, is_regular_pioneer, is_special_pioneer, spiritual_status
FROM publishers
WHERE id = $1;
