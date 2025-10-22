# Build the application
all: build test

build: build-css
	@echo "Building..."
	@go build -o main.exe cmd/wryte/main.go

# Build Tailwind CSS
build-css:
	@echo "Building Tailwind CSS..."
	@npm run build

# Run the application
run:
	@go run cmd/wryte/main.go

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main.exe

# Watch Tailwind CSS only
watch-css:
	@echo "Watching Tailwind CSS..."
	@npm run watch

watch:
	@powershell -ExecutionPolicy Bypass -Command "if (Get-Command air -ErrorAction SilentlyContinue) { \
		air; \
		Write-Output 'Watching...'; \
	} else { \
		Write-Output 'Installing air...'; \
		go install github.com/air-verse/air@latest; \
		air; \
		Write-Output 'Watching...'; \
	}"

# Docker commands
docker-build:
	@echo "Building Docker image..."
	@docker build -t wryte:latest .

docker-run:
	@echo "Running Docker container..."
	@docker run -p 8080:8080 wryte:latest

docker-up:
	@echo "Starting with docker-compose..."
	@docker-compose up -d

docker-down:
	@echo "Stopping docker-compose..."
	@docker-compose down

docker-logs:
	@echo "Showing logs..."
	@docker-compose logs -f

docker-rebuild:
	@echo "Rebuilding and restarting..."
	@docker-compose up -d --build

.PHONY: all build build-css run test clean watch watch-css docker-build docker-run docker-up docker-down docker-logs docker-rebuild
