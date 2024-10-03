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
RUN go build -o main .

# Use a smaller base image for the final container
FROM alpine:latest

# Set work directory
WORKDIR /root/

# Copy the Go app binary from the builder
COPY --from=builder /app/main .

# Command to run the app
CMD ["./server/main"]
