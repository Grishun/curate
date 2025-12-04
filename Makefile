.PHONY:

generate:
	@go tool oapi-codegen -config ./configs/oapi-config.yaml ./api/openapi.yaml

test:
	@go test ./...

docker-build:
	@docker build . -t curate-dev

run:
	@docker compose up -d

down:
	@docker compose down

purge:
	@docker compose down -v --remove-orphans