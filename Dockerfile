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

# Build binary
RUN go build -o main .

# ---- Run Stage ----
FROM alpine:3.19
WORKDIR /app

# Copy built binary
COPY --from=builder /app/main .

# Expose backend port
EXPOSE 8081

# Run the backend
CMD ["./main"]
