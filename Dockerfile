# Use a minimal Go image for building
FROM golang:1.23.5-alpine AS builder

# Set necessary environment variables for cross-compilation
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Set working directory
WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod tidy && go mod download

# Copy the rest of the application source code
COPY . .

# Build the application binary
RUN go build -o receipt-processor ./cmd/receipt-processor

# Use a minimal runtime image for the final container
FROM alpine:latest

# Set working directory
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/receipt-processor .

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./receipt-processor"]
