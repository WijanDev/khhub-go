# Implementation plan

1. Count failures per username/email and per IP in memory or a small table. After N failures, 429 for a cool-down.
2. Same message whether the user exists or not. Log with slog, no password.
3. Verify: httptest N+1 returns 429; success clears the counter; `go test ./...`.
