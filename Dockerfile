# Build stage
FROM golang:1.16-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY backend/go.mod backend/go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY backend/*.go ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Copy the frontend files
COPY frontend ./frontend

# Copy the test_image and test_result directories
COPY test_image ./test_image
COPY test_result ./test_result

# Expose port 8080
EXPOSE 8080

# Run the binary
CMD ["./main"]
