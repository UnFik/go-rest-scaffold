# Build stage
FROM golang:1.23-alpine AS builder

# Install required build tools
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/main ./cmd/web/main.go

# Runtime stage
FROM alpine:latest

# Install required runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user for security
RUN adduser -D -g '' appuser

# Set working directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy config files
COPY --from=builder /app/config.json .

# Copy migration files (optional, for running migrations)
COPY --from=builder /app/db/migrations ./db/migrations

# Change ownership to non-root user
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose the application port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:3000/api/health || exit 1

# Run the application
CMD ["./main"]
