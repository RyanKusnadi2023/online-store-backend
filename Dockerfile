# Use official Go image as the base image
FROM golang:1.23-alpine as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy everything from the current directory to /app inside the container
COPY . .

# Download dependencies (go.mod & go.sum)
RUN go mod tidy

# Build the Go app
RUN go build -o main cmd/server/main.go

# Start a new stage from scratch
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the Go app
CMD ["./main"]
