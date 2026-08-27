# Implementation plan

1. Table-driven `domain.ServiceYear` cases on the August/September boundary, plus `ServiceYearBounds` for 2025.
2. Table-driven `domain.ReportsHours` cases (publisher, RP/SP, auxiliary, missing/negative hours) in `backend/internal/domain/domain_test.go`. Remove the duplicate `TestServiceYear` from `service`.
3. Verify: `cd backend && go test ./internal/domain/ ./...`.
