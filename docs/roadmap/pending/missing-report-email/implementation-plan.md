# Implementation plan

1. SMTP config in env (not committed). Skip send when unset (same flag as invite email).
2. Secretary action: list missing (reuse `missing-report-list`) and send. One mail per publisher; no BCC of the whole list in one body.
3. Verify: no SMTP → button hidden or 503; fixture does not leak other congregations; `go test ./...` with a fake sender.
