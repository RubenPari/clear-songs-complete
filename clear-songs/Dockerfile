# Build stage
FROM golang:alpine AS build

WORKDIR /app

# Download dependencies first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o clear-songs ./cmd/server/main.go

# Production stage
FROM alpine:latest

WORKDIR /app

# Install CA certificates for HTTPS requests (e.g., Spotify API)
RUN apk --no-cache add ca-certificates

# Copy the binary from the build stage
COPY --from=build /app/clear-songs .

# Run the binary
CMD ["./clear-songs"]
