FROM golang:1.21 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make build

FROM node:22.12.0-alpine AS tailwind-builder
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm install
COPY tailwind.config.js ./tailwind.config.js
COPY ./internal/static/css ./internal/static/css
RUN npx tailwindcss build ./internal/static/css/main.css -o ./internal/static/css/tailwind.css --minify

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/bin/gohtmx ./gohtmx
COPY --from=tailwind-builder /app/internal/static ./internal/static
COPY ./internal/templates ./internal/templates
EXPOSE 8080
CMD ["./gohtmx"]
