# Implementation plan

1. `privilege_events` (publisher_id, kind, started_on, ended_on). On publisher flag change, close or open a row.
2. Show history on the secretary publisher form. Read-only for the person on their own card if useful.
3. Verify: toggling RP twice writes two intervals; `go test ./...`.
