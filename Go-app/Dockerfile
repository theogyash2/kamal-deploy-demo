# Stage 1: Build the Go binary
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o my-go-app .

# Stage 2: Create the minimal production image
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/my-go-app .

# Expose the web port (Kamal defaults to 80)
EXPOSE 3001

# Define the entry point. The CMD is what Kamal will override.
ENTRYPOINT ["/app/my-go-app"]