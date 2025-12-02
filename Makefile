.PHONY:

generate:
	@go tool oapi-codegen -config ./configs/oapi-config.yaml ./api/openapi.yaml

test:
	@go test ./...

build-docker:
	@docker build . -t curate-dev