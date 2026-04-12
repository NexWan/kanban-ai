include .env
export

infra-up:
	docker compose -f infra/docker/compose.infra.yml up -d

infra-down:
	docker compose -f infra/docker/compose.infra.yml down

api:
	cd apps/api && go run ./cmd/server

ai:
	cd apps/ai-service && uvicorn app.main:app --reload --port 8001

web:
	cd apps/web && npm run dev

migrate-up:
	migrate -path apps/api/internal/db/migrations \
	-database "$(DATABASE_URL)" \
	up

migrate-down:
	migrate -path apps/api/internal/db/migrations \
	-database "$(DATABASE_URL)" \
	down 1

migrate-create:
	migrate create -ext sql -dir apps/api/internal/db/migrations -seq $(name)