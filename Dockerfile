FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o gohtmx main.go

# Tailwind CSS builder stage
FROM node:22.12.0-alpine AS tailwind-builder
WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci --only=production

COPY tailwind.config.js ./
COPY ./internal/static/css ./internal/static/css
COPY ./internal/template ./internal/template

# Build Tailwind CSS
RUN npx tailwindcss build ./internal/static/css/main.css -o ./internal/static/css/tailwind.css --minify

# Final production stage
FROM scratch

# Copy ca-certificates for HTTPS requests
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Create non-root user
COPY --from=builder /etc/passwd /etc/passwd

# Copy application binary
COPY --from=builder /app/gohtmx /gohtmx

# Copy static assets and templates
COPY --from=tailwind-builder /app/internal/static /internal/static
COPY --from=builder /app/internal/template /internal/template

# Copy configuration files
COPY config.yaml /config.yaml

# Create non-root user for security
USER nobody

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD ["/gohtmx", "--health-check"] || exit 1

# Set environment to production
ENV GOHTMX_APP_ENVIRONMENT=production

# Run the application
CMD ["/gohtmx"]
