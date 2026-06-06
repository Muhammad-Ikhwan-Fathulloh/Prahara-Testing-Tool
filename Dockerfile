# Stage 1: Build Vue Frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /web

# Copy package files
COPY web/package.json ./

# Install dependencies with legacy peer deps flag for compatibility
RUN npm install --legacy-peer-deps

# Copy source code
COPY web/ ./

# Build the frontend
RUN npm run build

# Stage 2: Build Go Backend
FROM golang:1.21-alpine AS backend-builder
WORKDIR /app
ENV CGO_ENABLED=0

# Copy go module files
COPY api/go.mod ./
COPY api/go.sum* ./

# Download dependencies
RUN go mod download || true

# Copy API source code
COPY api/ ./

# Tidy and build
RUN go mod tidy && go build -o main .

# Stage 3: Production Runner
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates wget

# Install k6
RUN wget https://github.com/grafana/k6/releases/download/v0.45.0/k6-v0.45.0-linux-amd64.tar.gz \
    && tar -xvzf k6-v0.45.0-linux-amd64.tar.gz \
    && mv k6-v0.45.0-linux-amd64/k6 /usr/bin/k6 \
    && rm -rf k6-v0.45.0-linux-amd64*

WORKDIR /app

# Copy Go binary
COPY --from=backend-builder /app/main .

# Copy Frontend build
COPY --from=frontend-builder /web/dist ./web/dist

EXPOSE 3000

CMD ["./main"]
