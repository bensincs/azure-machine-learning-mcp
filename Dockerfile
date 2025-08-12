# Build stage
FROM golang:1.24.6-alpine AS builder

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mcp-server ./cmd/mcp-server

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests and Azure CLI dependencies
RUN apk --no-cache add ca-certificates curl bash python3 py3-pip

# Install Azure CLI
RUN pip3 install azure-cli --break-system-packages

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
  adduser -u 1001 -S appuser -G appgroup

# Create .azure directory for the user
RUN mkdir -p /home/appuser/.azure && \
  chown -R appuser:appgroup /home/appuser/.azure

WORKDIR /home/appuser/

# Copy the binary from builder stage
COPY --from=builder /app/mcp-server .

# Change ownership to non-root user
RUN chown appuser:appgroup /home/appuser/mcp-server

# Switch to non-root user
USER appuser

# Expose port (if your MCP server uses a specific port)
# EXPOSE 8080

# Command to run the application
CMD ["./mcp-server"]
