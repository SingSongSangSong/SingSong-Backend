# Use the official Golang image as a build stage
FROM golang:1.22 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .
# Copy the .env file into the Docker image
COPY .env .env
# Start a new stage from scratch
FROM golang:1.22

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source code from the builder stage
COPY --from=builder /app .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the Go application
CMD ["go", "run", "cmd/main.go"]