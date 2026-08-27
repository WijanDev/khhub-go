# Implementation plan

1. Add a gitleaks (or equivalent) step to `.github/workflows/ci.yml`.
2. Commit a baseline so existing false positives are reviewed once.
2. Confirm a dummy secret in a branch fails CI.
