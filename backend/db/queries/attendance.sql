-- name: ListAttendance :many
SELECT id, meeting_date, kind, in_person, online, created_at
FROM meeting_attendance
WHERE meeting_date >= sqlc.arg(from_date) AND meeting_date <= sqlc.arg(to_date)
ORDER BY meeting_date, kind;

-- name: UpsertAttendance :one
INSERT INTO meeting_attendance (meeting_date, kind, in_person, online)
VALUES ($1, $2, $3, $4)
ON CONFLICT (meeting_date, kind) DO UPDATE SET
    in_person = EXCLUDED.in_person,
    online = EXCLUDED.online
RETURNING id, meeting_date, kind, in_person, online, created_at;

-- name: DeleteAttendance :exec
DELETE FROM meeting_attendance WHERE id = $1;
