# ==========================================
# STAGE 1: Build the Go Application
# ==========================================
FROM golang:1.26 AS builder

# Set the current working directory inside the container
WORKDIR /app

# Copy the Go module manifests to cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of your application source code
COPY . .

# Build a statically linked Linux binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/notes-api ./cmd/server

# ==========================================
# STAGE 2: Create a lightweight runtime
# ==========================================
FROM alpine:latest

WORKDIR /app

# Install CA certificates for making secure HTTPS calls (common for Go apps)
RUN apk --no-cache add ca-certificates

# Copy the migration folder specifically
COPY --from=builder /app/internal/db ./internal/db

# Copy the go binary from the builder stage
COPY --from=builder /bin/notes-api /bin/notes-api

# Heroku conventionally uses the PORT environment variable
ENV PORT=8080
EXPOSE 8080

# Run the compiled binary
CMD ["/bin/notes-api"]