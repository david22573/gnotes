build:
	@go build -o bin/app_prod cmd/app/main.go

run:
	@bin/app_prod
