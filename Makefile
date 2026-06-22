up:
	docker compose -f deploy/docker/compose.yaml up -d

down:
	docker compose -f deploy/docker/compose.yaml down

reset:
	docker compose -f deploy/docker/compose.yaml down -v
	docker compose -f deploy/docker/compose.yaml up -d

run:
	set -a && source .env && set +a && go run ./cmd/api