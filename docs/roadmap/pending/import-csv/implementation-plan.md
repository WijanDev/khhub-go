# Implementation plan

1. Document column names. `POST` multipart, parse on the API, validate, insert in a transaction. Reject hours for non-pioneers per `domain.ReportsHours`.
2. Dry-run preview + commit. Spanish errors per row.
3. Verify: fixture CSV; bad hours 400; no cross-congregation writes; `go test ./...`.
