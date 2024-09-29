# First stage: Build the Go app
FROM golang:1.23.1-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download 

# Copy the source code
COPY . .

# Build the Go app with CGO disabled for a statically linked binary
RUN CGO_ENABLED=0 go build -o main .

# Second stage: Create a smaller image
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Expose the port
EXPOSE 1337

# Command to run the executable
CMD ["./main"]
