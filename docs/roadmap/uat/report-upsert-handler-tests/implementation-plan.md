# Implementation plan

1. Extract a `reportStore` interface (same idea as `sessionQuerier`) so `putReports` can be httptested without Postgres.
2. Table-driven `PUT /reports` tests around `domain.ReportsHours`: regular/special pioneer, auxiliary that month, publisher without hours, and hours rejected on a non-pioneer.
3. Run `go test ./...`.
