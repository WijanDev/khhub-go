# khhub-go skill and backend alignment

- **Slug:** khhub-go
- **Status:** proposed
- **Merge date:**
- **App version:**

## Summary

Add a project skill that tells agents how to write and review Go in this repo (package by concern, `http` / `service` / `domain` / `store`, no one-type-per-file, no pass-through services). In the same pass, move activity recompute into `service`, use `domain.IsHourReporter` in the report list, and stop sending raw Go errors on password change and hour validation. Secretaries only notice the last part: those failures stay in Spanish.

Design: [docs/superpowers/specs/2026-08-27-khhub-go-design.md](../../../superpowers/specs/2026-08-27-khhub-go-design.md).

## Out of scope

- Splitting `auth.go` or adding a service per resource.
- Renaming `http` / `store` to match generic Gin tutorials.
- A reports or publishers httptest suite (see `report-upsert-handler-tests`).
- Changing generated sqlc files or the frontend.
