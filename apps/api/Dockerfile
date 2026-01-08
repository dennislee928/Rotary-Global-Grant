FROM golang:1.22-alpine AS build

WORKDIR /src

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Debug: Log build environment and verify source files
RUN echo "=== Build Debug Info ===" && \
    echo "Working directory: $(pwd)" && \
    echo "Go version: $(go version)" && \
    echo "Module path: $(head -1 go.mod | cut -d' ' -f2)" && \
    echo "Source files in cmd/server:" && \
    ls -la ./cmd/server/ || echo "ERROR: cmd/server directory not found" && \
    echo "=== Starting Build ==="

# Build the binary with verbose output
# Fix: Use single quotes for ldflags to avoid shell interpretation issues
RUN CGO_ENABLED=0 GOOS=linux go build -v -x -ldflags='-w -s' -o /out/server ./cmd/server 2>&1 | tee /tmp/build.log || \
    (echo "=== Build Failed ===" && \
     cat /tmp/build.log && \
     echo "=== Go Environment ===" && \
     go env && \
     echo "=== Module Graph ===" && \
     go mod graph | head -20 && \
     exit 1)

# Final stage
FROM gcr.io/distroless/static:nonroot

COPY --from=build /out/server /server

USER nonroot:nonroot

ENV HTTP_ADDR=:8080

EXPOSE 8080

ENTRYPOINT ["/server"]
