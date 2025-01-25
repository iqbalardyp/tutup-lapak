# Use the official Golang base image with version 1.21.4
FROM golang:1.23.4-alpine as builder

# Set the working directory inside the container
WORKDIR /app

# Copy dependency files first to leverage Docker caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN GOARCH=amd64 GOOS=linux go build -o main ./cmd/api

# Use a minimal base image for the final image
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the built executable from the builder stage
COPY --from=builder /app/main .

# Expose the port the application will run on
EXPOSE 8080

# Command to create the .env file and run the application
CMD sh -c 'printenv > .env && ./main'