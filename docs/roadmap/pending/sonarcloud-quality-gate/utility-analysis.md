# Utility analysis

- **Score:** 4
- **Scale:** 1 = low … 5 = high

## Justification

Stops “we will add tests later” on new HTTP and UI. The secretary never sees this; the operator (and agents) do, on every PR. Skipping it leaves only `go test ./...` with no coverage floor. That is production-safety for a public repo, not a weekly congregation task, but it unblocks trusting `dev` and `main`.
