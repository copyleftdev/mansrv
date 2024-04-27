# Use a base image with Go and a full OS to include man pages
FROM golang:1.18-buster

# Install man, groff, gzip, and man pages
RUN apt-get update && \
    apt-get install -y man manpages manpages-dev man-db groff gzip

# Set the working directory in the container
WORKDIR /app

# Copy the Go modules and sum files
COPY go.mod go.sum ./

# Download Go dependencies
RUN go mod download

# Copy only the relevant Go source code file into the container
COPY server.go .

# Build the specific application
RUN go build -o server ./server.go

# Expose port 8080 to the outside world
EXPOSE 8887

# Command to run the executable
CMD ["./server"]
