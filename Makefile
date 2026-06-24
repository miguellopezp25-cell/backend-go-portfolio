.PHONY: build run test migrate sqlc docker-build docker-up clean vet lint help

APP_NAME = app

build:
	go build -o bin/$(APP_NAME) ./main.go

run:
	go run ./main.go serve --config config.yaml

test:
	go test ./... -v -count=1

migrate:
	go run ./main.go migrate --config config.yaml

sqlc:
	sqlc generate

docker-build:
	docker compose build

docker-up:
	docker compose up -d

vet:
	go vet ./...

lint:
	golangci-lint run ./...

clean:
	rm -rf bin/

help:
	@echo "Targets:"
	@echo "  build       - Compile the project"
	@echo "  run         - Start the HTTP server"
	@echo "  test        - Run all tests"
	@echo "  migrate     - Run database migrations"
	@echo "  sqlc        - Regenerate SQLC code"
	@echo "  docker-build - Build Docker images"
	@echo "  docker-up   - Start services with Docker Compose"
	@echo "  vet         - Run go vet"
	@echo "  lint        - Run golangci-lint"
	@echo "  clean       - Remove build artifacts"
