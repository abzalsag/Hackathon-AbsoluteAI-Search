# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o main cmd/api/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary
COPY --from=builder /app/main .
COPY --from=builder /app/db ./db
COPY --from=builder /app/configs ./configs

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]