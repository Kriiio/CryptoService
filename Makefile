# ==============================================================================
# Main
run:
	go run ./cmd//main.go

build:
	go build ./cmd//main.go

test:
	go test ./...
# ==============================================================================
# Tools commands
lint:
	echo "Starting linters"
	golangci-lint run 
# ==============================================================================
# Docker commands
docker-build:
	docker-compose build --no-cache

docker-up:
	docker-compose up -d
# ==============================================================================
# Цели
.PHONY: run build test lint docker-build docker-up
	