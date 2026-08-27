# Implementation plan

1. Congregation or publisher setting for annual/monthly hour goal (defaults: current RP/SP expectations as data, not hardcoded theology in many places).
2. `GET` own progress from summed `field_service_reports` hours. 403 if not RP/SP (unless secretary).
3. SPA page. Verify: another pioneer’s id is 403; `go test ./...`.
