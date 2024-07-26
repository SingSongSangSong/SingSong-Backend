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

# Accept build arguments and set environment variables
ARG PINECONE_API_KEY
ARG PINECONE_INDEX
ARG DB_HOST
ARG DB_PORT
ARG DB_NAME
ARG DB_USER
ARG DB_PASSWORD
ARG REDIS_ADDR
ARG REDIS_PASSWORD
ARG SECRET_KEY
ARG KAKAO_REST_API_KEY
ARG KAKAO_ISSUER
ARG JWT_ISSUER
ARG JWT_ACCESS_VALIDITY_SECONDS
ARG JWT_REFRESH_VALIDITY_SECONDS

ENV PINECONE_API_KEY=$PINECONE_API_KEY
ENV PINECONE_INDEX=$PINECONE_INDEX
ENV DB_HOST=$DB_HOST
ENV DB_PORT=$DB_PORT
ENV DB_NAME=$DB_NAME
ENV DB_USER=$DB_USER
ENV DB_PASSWORD=$DB_PASSWORD
ENV REDIS_ADDR=$REDIS_ADDR
ENV REDIS_PASSWORD=$REDIS_PASSWORD
ENV SECRET_KEY=$SECRET_KEY
ENV KAKAO_REST_API_KEY=$KAKAO_REST_API_KEY
ENV KAKAO_ISSUER=$KAKAO_ISSUER
ENV JWT_ISSUER=$JWT_ISSUER
ENV JWT_ACCESS_VALIDITY_SECONDS=$JWT_ACCESS_VALIDITY_SECONDS
ENV JWT_REFRESH_VALIDITY_SECONDS=$JWT_REFRESH_VALIDITY_SECONDS

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the Go application
CMD ["./main"]