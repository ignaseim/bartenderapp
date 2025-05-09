# Build stage
FROM golang:1.22-alpine AS builder

# Install necessary build tools
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY auth/go.mod auth/go.mod
COPY pkg/go.mod pkg/go.mod

# Copy the source code
COPY auth/ auth/
COPY pkg/ pkg/

# Build the application
WORKDIR /app/auth
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o auth-service ./cmd/main.go

# Final stage
FROM alpine:latest

# Install dependencies
RUN apk --no-cache add ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/auth/auth-service .

# Expose port
EXPOSE 8081

# Run the application
CMD ["./auth-service"] 