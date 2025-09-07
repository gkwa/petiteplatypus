FROM golang:1.23-alpine AS builder

# Install git for go modules
RUN apk add --no-cache git

# Set GOTOOLCHAIN to auto to allow downloading newer Go versions if needed
ENV GOTOOLCHAIN=auto

# Set working directory
WORKDIR /build

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./

# Download dependencies (this layer will be cached)
RUN go mod download

# Copy only the Go source files, not the run.sh script
COPY main.go ./
COPY templates/ ./templates/

# Install the Go tools and build our app in a single RUN instruction
RUN go install github.com/gkwa/petiteplatypus@latest && \
    go install github.com/Yakitrak/obsidian-cli@latest && \
    go build -o petiteplatypus .

# Use alpine for the final image
FROM alpine:3.20

# Install bash for script execution
RUN apk add --no-cache bash

# Copy the Go binaries from builder stage (this won't change often)
COPY --from=builder /go/bin/obsidian-cli /usr/local/bin/
COPY --from=builder /build/petiteplatypus /usr/local/bin/

# Create working directory
WORKDIR /app

# Copy the script file (this is the only layer that changes)
COPY run.sh /app/run.sh
RUN chmod +x /app/run.sh

# Default command
CMD ["/app/run.sh"]
