# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Ensure we use pure Go (no CGO) for maximum compatibility
ENV CGO_ENABLED=0

# Copy go mod files
COPY go.mod ./
COPY go.sum* ./

# Copy the source code
COPY . .

# Ensure go.sum is present and correct
RUN go mod tidy

# Build the application
RUN go build -o main .

# Final stage
FROM alpine:latest

# Install dependencies for k6
RUN apk add --no-cache ca-certificates wget

# Install k6
RUN wget https://github.com/grafana/k6/releases/download/v0.45.0/k6-v0.45.0-linux-amd64.tar.gz \
    && tar -xvzf k6-v0.45.0-linux-amd64.tar.gz \
    && mv k6-v0.45.0-linux-amd64/k6 /usr/bin/k6 \
    && rm -rf k6-v0.45.0-linux-amd64*

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .

EXPOSE 3000

CMD ["./main"]
