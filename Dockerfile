FROM golang:1.25-alpine AS build

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache curl git ca-certificates libstdc++ libgcc

# Download Go dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build Tailwind CSS using standalone CLI
RUN curl -sL https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64-musl -o tailwindcss && \
    chmod +x tailwindcss && \
    ./tailwindcss -i web/styles/tailwind.css -o web/assets/css/output.css --minify

# Build Go application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main cmd/wryte/main.go

FROM alpine:latest AS prod

WORKDIR /app

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Copy binary from build stage
COPY --from=build /app/main /app/main

# Set environment variables
ENV HOST=0.0.0.0
ENV PORT=8080
ENV ENV=production

EXPOSE 8080

CMD ["./main"]
