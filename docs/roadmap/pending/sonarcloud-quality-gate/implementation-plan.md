# Implementation plan

1. Create the SonarCloud project for `WijanDev/khhub-go` (or the current GitHub repo). Bind GitHub. Add `SONAR_TOKEN` as a GitHub Actions secret. Do not commit the token.
2. Add `sonar-project.properties` at the repo root: project key, sources `backend` + `frontend/src`, exclusions for `backend/internal/store/*.sql.go`, `frontend/dist`, `node_modules`.
3. In `.github/workflows/ci.yml`, after tests: run SonarCloud scan with coverage reports. Go: `go test ./... -coverprofile=coverage.out` in `backend`. Frontend: enable coverage when a test runner exists; until then scan without TS coverage rather than inventing a fake report.
4. Write `docs/quality.md` (English) with the table from `idea.md`. One line in `AGENTS.md` and `README.md` pointing at it. PR template (if none): “quality bar in `docs/quality.md`”.
5. Set the SonarCloud quality gate to the new-code conditions in `idea.md`. Require the Sonar check on PRs into `dev` (and `main` until `dev-branch-workflow` lands).
6. Verify: a PR that only adds untested Go in `backend/internal/http` fails the gate. A docs-only PR passes. `go test ./...` and frontend lint/build still run first.
