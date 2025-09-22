FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod file
COPY go.mod ./

# Download dependencies (this will create go.sum)
RUN go mod download

# Copy source code
COPY . .

# Tidy up dependencies and generate go.sum
RUN go mod tidy

# Build the application
RUN go build -o main .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Expose ports
EXPOSE 8084
EXPOSE 3013

# Run the application
CMD ["./main"]