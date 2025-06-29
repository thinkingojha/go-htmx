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
RUN CGO_ENABLED=0 GOOS=linux go build -o gohtmx main.go

# Tailwind CSS builder stage
FROM node:22.12.0-alpine AS tailwind-builder
WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY tailwind.config.js ./
COPY ./internal/static/css ./internal/static/css
COPY ./internal/template ./internal/template

# Build Tailwind CSS
RUN npm run build-css

# Final production stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy application binary
COPY --from=builder /app/gohtmx /gohtmx

# Copy static assets and templates
COPY --from=tailwind-builder /app/internal/static /internal/static
COPY --from=builder /app/internal/template /internal/template

# Copy configuration files and blog content
COPY config.production.yaml /config
COPY blogs /blogs
COPY experience.yaml /experience.yaml

# Create non-root user for security
RUN adduser -D -s /bin/sh appuser
USER appuser

# Expose port
EXPOSE 8080

# Run the application
CMD ["/gohtmx"]
