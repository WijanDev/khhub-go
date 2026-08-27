# Implementation plan

1. Table or congregation-scoped row: `(congregation_id, year, month, locked_at, locked_by)`.
2. Report and attendance writes return 409 when locked unless secretary unlocks. Publishers cannot unlock.
3. SPA: lock/unlock on the dashboard or reports header. Verify: httptest lock then upsert; `go test ./...`.
