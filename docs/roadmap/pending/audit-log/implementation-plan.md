# Implementation plan

1. `audit_events` (congregation_id, actor_user_id, action, entity, entity_id, at). Write from handlers after successful mutations.
2. Secretary/superadmin list page, filtered, no PII dumps beyond names already on screen.
3. Verify: publisher 403; a report upsert creates a row; `go test ./...`.
