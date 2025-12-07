.PHONY:

oapi-generate:
	@go tool oapi-codegen -config ./configs/oapi-config.yaml ./api/openapi.yaml

protoc:
	@protoc \
	    --proto_path=./api \
	    --go_out=./internal/transport/grpc/generated --go_opt=paths=source_relative \
	    --go-grpc_out=./internal/transport/grpc/generated --go-grpc_opt=paths=source_relative \
	    ./api/*.proto

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