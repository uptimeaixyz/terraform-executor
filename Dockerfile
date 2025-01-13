# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

# Final stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
RUN chmod +x /app/main

# Run the binary
CMD ["./main"]