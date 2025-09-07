FROM golang:1.23-alpine AS builder

# Install git for go modules
RUN apk add --no-cache git

# Set GOTOOLCHAIN to auto to allow downloading newer Go versions if needed
ENV GOTOOLCHAIN=auto

# Set working directory
WORKDIR /build

# Copy source code and templates
COPY . .

# Install the Go tools and build our app in a single RUN instruction
RUN go install github.com/gkwa/petiteplatypus@latest && \
    go install github.com/Yakitrak/obsidian-cli@latest && \
    go build -o petiteplatypus .

# Use alpine for the final image
FROM alpine:3.20

# Install bash for script execution
RUN apk add --no-cache bash

# Copy the Go binaries from builder stage
COPY --from=builder /go/bin/obsidian-cli /usr/local/bin/
COPY --from=builder /build/petiteplatypus /usr/local/bin/

# Create working directory
WORKDIR /app

# Create the commands as a script
RUN echo '#!/bin/bash' > /app/run.sh && \
    echo 'petiteplatypus generate /tmp/trash' >> /app/run.sh && \
    echo 'obsidian-cli set-default /tmp/trash' >> /app/run.sh && \
    echo 'obsidian-cli print-default' >> /app/run.sh && \
    chmod +x /app/run.sh

# Default command
CMD ["/app/run.sh"]
