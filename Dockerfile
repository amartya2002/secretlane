############################
# Build stage
############################
FROM golang:1.23-bullseye AS builder

WORKDIR /src

# Enable Go modules and turn on CGO (needed for sqlite3)
ENV CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

# Download dependencies first (better layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Build the binary
RUN go build -o /bin/secretlane ./...

############################
# Runtime stage (distroless)
############################
FROM gcr.io/distroless/base-debian12

# Working directory for the app
WORKDIR /app

# Copy binary and config
COPY --from=builder /bin/secretlane /app/secretlane
COPY config.yaml /app/config.yaml

# Expose default port (overridable via PORT env)
EXPOSE 8080

# Non-root user for better security (UID 65532 is 'nonroot' in distroless)
USER 65532:65532

# Entrypoint
ENTRYPOINT ["/app/secretlane"]
