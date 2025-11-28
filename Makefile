.PHONY:

generate:
	@go tool oapi-codegen -config ./configs/oapi-config.yaml ./api/openapi.yaml