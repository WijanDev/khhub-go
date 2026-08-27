# Implementation plan

1. Add `httptest` cases in `backend/internal/http` for login success, bad password, missing cookie, and `/auth/me`.
2. Assert `Set-Cookie` name `khhub_session` and 401 JSON on protected routes.
2. Run `go test ./...`.
