# Implementation plan

1. Add `GET /publishers/:id/service-year?year=` (secretary; publisher may read **own** id only). Return 12 months from existing `field_service_reports`.
2. SPA route on the publisher record (and “my card” for a logged-in publisher). Spanish labels. Depend on `accounts-multicongregation` for `/auth/me`.
3. Verify: httptest own vs other publisher (403); `go test ./...`; open the card in the browser for a seed pioneer and a publisher without hours.
