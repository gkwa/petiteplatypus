FROM golang:1.23-alpine AS builder

# Install git for go modules
RUN apk add --no-cache git

# Set GOTOOLCHAIN to auto to allow downloading newer Go versions if needed
ENV GOTOOLCHAIN=auto

# Install the Go tools
RUN go install github.com/gkwa/petiteplatypus@latest
RUN go install github.com/Yakitrak/obsidian-cli@latest

# Use alpine for the final image
FROM alpine:latest

# Install bash for script execution
RUN apk add --no-cache bash

# Copy the Go binaries from builder stage
COPY --from=builder /go/bin/petiteplatypus /usr/local/bin/
COPY --from=builder /go/bin/obsidian-cli /usr/local/bin/

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
