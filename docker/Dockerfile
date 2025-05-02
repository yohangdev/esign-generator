# syntax=docker/dockerfile:1

# Use the official Golang image to build the app
FROM golang:1.24 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY src/go.mod src/go.sum ./
RUN go mod download

# Copy source code
COPY src/ .

# Build the Go app
RUN go build -o image-service

# Use a minimal runtime image
FROM debian:bookworm-slim

WORKDIR /app

# Copy the binary and the Arial font file
COPY --from=builder /app/image-service ./
COPY src/arial.ttf ./

# Expose the service port
EXPOSE 8080

# Run the service
CMD ["./image-service"]
