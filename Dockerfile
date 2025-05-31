# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go binary (adjust path if your main.go is elsewhere)
RUN CGO_ENABLED=0 GOOS=linux go build -o swagtask ./cmd/swagtask

# Final image
FROM alpine:latest
WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates curl


# Download migrate CLI binary
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz \
  && mv migrate.linux-amd64 /usr/local/bin/migrate

  # Copy the Go binary and migrations
COPY --from=builder /app/swagtask .
COPY internal/db/migrations/ internal/db/migrations/


# Copy static files and templates
COPY web/ web/
COPY internal/db/migrations/ internal/db/migrations/

# Expose the port your app listens on
EXPOSE 8080

# Set environment variable for production
ENV PORT=8080

# Start the app
CMD ["./swagtask"]
