# Use Golang base image
FROM alpinelinux/golang AS builder

# Set work directory inside the container
WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application files
COPY . .

# Build the Go app
RUN go build -o main ./cmd/server

# Use a smaller base image for the final container
FROM alpine:latest

# Set work directory
WORKDIR /simplevault

# Copy the Go app binary from the builder
COPY --from=builder /app/main .

# Copy SSL certs
COPY --from=builder /app/cert.pem .
COPY --from=builder /app/privatekey.pem .

# Command to run the app
CMD ["./main"]

# setting default for environment
ENV GK_ADDRESS=host.docker.internal:8080 \
    GK_CERT=cert.pem \
    GK_PRIVATE_KEY=privatekey.pem \
    GK_DATABASE="postgres://postgres:postgres@db:5432/postgres?sslmode=disable"
