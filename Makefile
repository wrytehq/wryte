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

# Database migration commands
migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Error: name parameter is required"; \
		echo "Usage: make migrate-create name=your_migration_name"; \
		exit 1; \
	fi
	@bash -c ' \
		latest=$$(ls internal/database/migrations/*.up.sql 2>/dev/null | sed "s/.*\/\([0-9]*\)_.*/\1/" | sort -n | tail -1); \
		if [ -z "$$latest" ]; then \
			next=000001; \
		else \
			next=$$(printf "%06d" $$((10#$$latest + 1))); \
		fi; \
		touch "internal/database/migrations/$${next}_$(name).up.sql"; \
		touch "internal/database/migrations/$${next}_$(name).down.sql"; \
		echo "Created migration files:"; \
		echo "  internal/database/migrations/$${next}_$(name).up.sql"; \
		echo "  internal/database/migrations/$${next}_$(name).down.sql" \
	'

migrate-up:
	@echo "Running migrations..."
	@go run cmd/wryte/main.go

migrate-down:
	@echo "Rolling back last migration..."
	@echo "Not implemented yet - use migration rollback carefully"


.PHONY: all build build-css run test clean watch docker-build docker-run docker-up docker-down docker-logs docker-rebuild migrate-create migrate-up migrate-down
