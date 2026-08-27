# Implementation plan

1. Migration: `field_service_groups` (`congregation_id`, name, overseer `publisher_id`) and `publishers.group_id`. sqlc + seed groups.
2. Secretary CRUD. Narrow list DTO for publisher portal (name, family, household name, group, privileges).
3. Verify: isolation by congregation; publisher cannot see hours; `go test ./...` and browser on Hermanos / Grupos.
