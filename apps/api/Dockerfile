FROM golang:1.22-alpine AS build

WORKDIR /src

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /out/server ./cmd/server

# Final stage
FROM gcr.io/distroless/static:nonroot

COPY --from=build /out/server /server

USER nonroot:nonroot

ENV HTTP_ADDR=:8080

EXPOSE 8080

ENTRYPOINT ["/server"]
