# Use the Go 1.20 Alpine image as the base image
FROM golang:1.20-alpine

RUN go version
ENV GOPATH=/

# Copy the Go project files to the container's working directory
COPY ./ ./


RUN go mod download

# Build the Go binary
RUN go build -o app ./cmd/main.go

# Set the entry point for the container
CMD ["./app"]