.PHONY: build run dev test clean docker-build docker-run

build:
	@echo "Building..."
	go build -o bin/server ./cmd/server

run: build
	@echo "Starting server..."
	./bin/server

dev:
	@echo "Starting in development mode..."
	go run ./cmd/server

test:
	@echo "Running tests..."
	go test ./...

clean:
	@echo "Cleaning..."
	rm -rf bin/
	go clean

docker-build:
	docker build -t ecommerce-backend .

docker-run:
	docker run -p 8080:8080 --env-file .env ecommerce-backend