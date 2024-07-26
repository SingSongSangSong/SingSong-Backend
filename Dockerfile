# Use the official Golang image as a build stage
FROM golang:1.22 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to the container
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the rest of the source code into the container
COPY . .

# Build the Go application
RUN go build -o main .

# Use a minimal Docker image for the final stage
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the Go application
CMD ["./main"]
