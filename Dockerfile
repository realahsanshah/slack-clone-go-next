# Build stage
FROM golang:1.24-alpine AS builder

# Install ca-certificates and tools for SSL and database
RUN apk --no-cache add ca-certificates git

WORKDIR /app

# Install sqlc for code generation
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code and SQL files
COPY . .

# Generate database code with sqlc
RUN sqlc generate

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo -o main ./cmd

# Final stage
FROM scratch

# Copy ca-certificates for SSL
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy binary from builder
COPY --from=builder /app/main /main

# Expose port
EXPOSE 8080

# Add health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD ["/main", "--health-check"] || exit 1

# Run the application
ENTRYPOINT ["/main"] 