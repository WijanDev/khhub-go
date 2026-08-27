# Implementation plan

1. Tables for talks (date, outline number/title, speaker name, congregation of speaker) scoped by `congregation_id`.
2. Secretary CRUD. Publisher portal: read the upcoming program (same permission as “program”).
3. Verify: tenant isolation; student can read program; `go test ./...` and browser.
