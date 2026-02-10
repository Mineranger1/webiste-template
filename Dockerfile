# Build Stage
FROM golang:1.24-alpine AS builder

# Install build dependencies for CGO (sqlite3)
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o main cmd/server/main.go

# Final Stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS calls (if any) and sqlite libs
RUN apk add --no-cache ca-certificates sqlite-libs

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy templates and static files
# Preserving directory structure relative to execution
COPY --from=builder /app/web ./web

# Set environment variables
ENV DB_PATH=biomix.db
ENV PORT=8080

# Expose the port
EXPOSE 8080

# Run the binary
CMD ["./main"]
