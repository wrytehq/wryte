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

.PHONY: all build build-css run test clean watch watch-css
