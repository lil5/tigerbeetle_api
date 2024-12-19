FROM golang:1.23.4-alpine AS builder

# Install git and build dependencies required by TigerBeetle
RUN apk update && apk add --no-cache git build-base zig
# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./
# Download dependencies
RUN --mount=type=cache,target=/go/pkg/mod \
   --mount=type=cache,target=/root/.cache/go-build \
   go mod download

# Copy the source code
COPY . .
# Enable CGO and specify the Zig compiler
ENV CGO_ENABLED=1 CC="zig cc"
# Build the Go application
RUN --mount=type=cache,target=/go/pkg/mod \
   --mount=type=cache,target=/root/.cache/go-build \
   go build -o tigerbeetle_api .


# Use a minimal base image for the final stage
FROM alpine:latest
# Set working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/tigerbeetle_api .

# Expose port (if your application listens on a port)
# EXPOSE 8080

# Set the entrypoint
ENTRYPOINT ["./tigerbeetle_api"]
