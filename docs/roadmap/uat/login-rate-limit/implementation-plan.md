# Implementation plan

1. In-memory limiter: count **failed** logins per IP and per normalized email. After 10 failures in 15 minutes on either key, respond 429 for the rest of the window. Success clears both keys. No schema (single API replica).
2. Same 401 body whether the user exists or not. Log lockouts and failures with `log/slog` (ip, email; never the password). Future passkey verify should reuse the same limiter.
3. Verify: httptest N+1 returns 429; a successful login clears the counter so N more failures are allowed; `go test ./...`.
