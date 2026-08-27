# Implementation plan

1. Use `gofmt -l` and `go vet ./...` in `backend/` (no golangci-lint config; avoids a lint-rule bikeshed).
2. Fail CI on those checks in `.github/workflows/ci.yml`. Local command: `make lint`.
3. Keep `*.go` as LF via `.gitattributes` so Windows `core.autocrlf` does not make `gofmt -l` fail locally.
4. Fix first-run `go vet` findings (keyed `MonthShare` literals in `activity_test.go`).
