# Stage 1: Build Vue Frontend
FROM node:18-alpine AS frontend-builder
WORKDIR /web
COPY web/package*.json ./
RUN npm install
COPY web/ ./
RUN npm run build

# Stage 2: Build Go Backend
FROM golang:1.20-alpine AS backend-builder
WORKDIR /app
ENV CGO_ENABLED=0
COPY api/go.mod api/go.sum* ./
RUN go mod download
COPY api/ ./
# We need the dist folder for the static file serving logic in Go (if checked at build time)
# but actually our Go code just references the path.
RUN go build -o main .

# Stage 3: Runner
FROM alpine:latest
RUN apk add --no-cache ca-certificates wget

# Install k6
RUN wget https://github.com/grafana/k6/releases/download/v0.45.0/k6-v0.45.0-linux-amd64.tar.gz \
    && tar -xvzf k6-v0.45.0-linux-amd64.tar.gz \
    && mv k6-v0.45.0-linux-amd64/k6 /usr/bin/k6 \
    && rm -rf k6-v0.45.0-linux-amd64*

WORKDIR /root/

# Copy Go binary
COPY --from=backend-builder /app/main .

# Copy Frontend build to the expected location for the Go server
COPY --from=frontend-builder /web/dist ./web/dist

EXPOSE 3000

CMD ["./main"]
