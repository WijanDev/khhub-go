.PHONY: dev-db api web test

dev-db:
	docker compose -f docker-compose.dev.yml --env-file .env up -d

api:
	cd backend && go run ./cmd/api

web:
	cd frontend && npm run dev

test:
	cd backend && go test ./internal/service/...
