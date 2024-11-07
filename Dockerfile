# Builder stage
FROM golang:1.22-alpine as builder
RUN apk update && apk add --no-cache git ca-certificates upx

WORKDIR /usr/src/app
ENV GO111MODULE=on
ENV GOPROXY=https://proxy.golang.org,direct

# Copy go.mod and go.sum for dependency resolution
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

# Copy the rest of the application source code
COPY . .

# Enable Go modules and build the application
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -ldflags="-s -w" -o bin/main main.go

# Optional: Compress the binary with UPX
RUN upx --best --lzma bin/main

# Executable image stage
FROM scratch

# Copy certificates and user information
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

# Copy the compiled Go binary
COPY --from=builder /usr/src/app/bin/main ./main

# Set the user (non-root user with ID)
USER 1000

# Expose the application's port
EXPOSE 8080

# Run the compiled binary
ENTRYPOINT ["./main"]
