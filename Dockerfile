# Build stage
FROM golang:1.23 AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o codema-server .

# Final stage
FROM gcr.io/distroless/static-debian12

WORKDIR /

# Copy the binary from the builder stage
COPY --from=builder /app/codema-server .

# Copy any static assets or templates
COPY --from=builder /app/static ./static

# Expose the port the app runs on
EXPOSE 8090

# Run the binary
CMD ["/codema-server"]
