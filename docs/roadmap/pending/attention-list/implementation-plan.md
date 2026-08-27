# Implementation plan

1. Derive the list from existing activity status (do not persist status as source of truth). Optional table for visit planned/done only (`publisher_id`, dates, no text).
2. `GET` for elder, MS, secretary, superadmin in the active congregation. 403 for publishers. Narrow DTO (directory fields + activity + visit flags).
3. SPA page. Verify: httptest roles; no remarks column; `go test ./...` and browser.
