build:
	@go build -o bin/app_prod cmd/app/main.go

run:
	@go run bin/app_prod
