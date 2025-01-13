# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

# Final stage
FROM hashicorp/terraform:1.10
WORKDIR /app
COPY --from=builder /app/main .
RUN chmod +x /app/main

ENV LISTEN 0.0.0.0:50051

ENTRYPOINT ["/app/main"]
