all: swag gqlgen run

run:
	@DEBUG=1 go run cmd/marketplace/main.go

swag:
	@~/go/bin/swag init --parseInternal --parseDependency -g cmd/marketplace/main.go

gqlgen:
	@go run github.com/99designs/gqlgen
