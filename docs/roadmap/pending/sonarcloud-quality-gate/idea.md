# SonarCloud and minimum quality bar

- **Slug:** sonarcloud-quality-gate
- **Status:** proposed
- **Merge date:**
- **App version:**

## Summary

Link the public repo to SonarCloud. CI fails the PR if the quality gate fails. Write the **minimum acceptable bar** in-repo so agents and humans use the same numbers.

## Proposed bar (new code on the PR)

These are the first-pass gates. Change them here if they are too strict or too weak; do not invent a second list in chat.

| Check | Minimum |
| --- | --- |
| `go test ./...` | Pass |
| `frontend`: `npm run lint` and `npm run build` | Pass |
| Sonar **new** bugs / vulnerabilities | 0 |
| Sonar **new** blocker / critical issues | 0 |
| Coverage on **new** Go and TS (excluding generated sqlc, `*.sql.go`, seed fixtures if noisy) | ≥ 90% |
| Duplicated lines on **new** code | ≤ 3% |
| Overall project coverage | Deferred. No fail-the-repo gate on day one. Configure a global % later. |

golangci-lint stays the separate idea `golangci-lint-ci`; when it exists, it is also required.

## Out of scope

- SonarQube on the VPS.
- Failing `main` on historical coverage.
- Paying for SonarCloud features we do not need (the public repo can use the free org plan).
