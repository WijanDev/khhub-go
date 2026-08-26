# khhub

Herramienta personal para la operación de una congregación: publicadores, informes mensuales de predicación y asistencia. **No es software oficial**, no habla con jw.org y no sustituye las tarjetas S-21 ni JW Hub.

Interfaz en español. Código y API en inglés.

## Requisitos locales

- Go 1.24+
- Node 22+
- Docker (solo para Postgres en desarrollo)

## Arranque

```bash
cp .env.example .env
# edita ADMIN_PASSWORD
docker compose -f docker-compose.dev.yml --env-file .env up -d
cd backend && go run ./cmd/api   # loads ../.env automatically
```

En otra terminal:

```bash
cd frontend && npm install && npm run dev
```

Abre http://localhost:5173 y entra con `ADMIN_EMAIL` / `ADMIN_PASSWORD`.

Recarga del API: `go install github.com/air-verse/air@latest` y `air` dentro de `backend/`.

## Año de servicio

1 de septiembre – 31 de agosto. Los informes siguen la práctica actual: todos marcan si participaron y los estudios; solo los precursores (y PA del mes) registran horas.

## Despliegue

Ver [docs/deploy-dokploy.md](docs/deploy-dokploy.md).
