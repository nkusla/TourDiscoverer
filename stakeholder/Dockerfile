FROM golang:1.23.4-alpine AS builder

WORKDIR /app

# Copy source code
COPY . .

# Download dependencies and generate go.sum
RUN go mod tidy
RUN go mod download

# Build the application
RUN go build -o main .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Expose port
EXPOSE 8081

# Run the binary
CMD ["./main"]
