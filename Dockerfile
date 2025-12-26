# ---- Build Stage ----
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Install git for Go modules
RUN apk add --no-cache git

# Copy Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary with memory optimizations
# -p 1 limits parallel compilation to reduce memory usage
RUN go build -p 1 -ldflags="-s -w" -o main .

# ---- Run Stage ----
FROM alpine:3.19
WORKDIR /app

# Create a non-root user for OpenShift
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Copy built binary
COPY --from=builder /app/main .

# Expose backend port
EXPOSE 8080

# Run the backend
CMD ["./main"]