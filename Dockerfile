# Stage 1: Build the Go application
FROM golang:1.22.5 AS builder

WORKDIR /app

# Copy the go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o ordersystem ./cmd/ordersystem

# Install the migrate CLI for handling database migrations
RUN go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Stage 2: Create the final image to run the application
FROM golang:1.22.5

WORKDIR /app

# Install netcat
RUN apt-get update && apt-get install -y netcat-openbsd

# Copy built application and migration tool from the builder image
COPY --from=builder /app/ordersystem .
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
COPY wait-for-mysql.sh /app/wait-for-mysql.sh

# Ensure the wait script is executable
RUN chmod +x /app/wait-for-mysql.sh

# Copy .env and migrations directory
COPY .env /app/.env
COPY ./internal/infra/database/migrations /app/migrations

# Set environment variables for MySQL
ENV DB_DRIVER=mysql
ENV DB_USER=root
ENV DB_PASSWORD=root
ENV DB_HOST=mysql
ENV DB_PORT=3306
ENV DB_NAME=orders

# Run migration command and start the application
ENTRYPOINT ["/bin/bash", "-c", "/app/wait-for-mysql.sh && migrate -path /app/migrations -database \"mysql://${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}\" up && ./ordersystem"]
