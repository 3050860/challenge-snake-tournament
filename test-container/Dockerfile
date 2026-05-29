# Test Container for Snake Tournament

FROM golang:1.23-alpine

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Copy the source
COPY . .

# Build the application
RUN go build -o snake-tournament .

# Install test dependencies
RUN apk add --no-cache bash

# Create a test user
RUN adduser -D -s /bin/sh testuser

# Set ownership
RUN chown -R testuser:testuser /app
USER testuser

# Expose port
EXPOSE 8080

# Command to run the application
CMD ["/app/snake-tournament"]