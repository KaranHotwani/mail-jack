# Stage 1: Builder
FROM golang:1.24-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to leverage Docker's layer caching
COPY go.mod .
COPY go.sum .

# Download Go modules
ENV GOTOOLCHAIN=auto
RUN go mod download

# Install runtime assets we will copy into the final image (CA certs, tzdata)
RUN apk add --no-cache ca-certificates tzdata

# Copy the rest of the application source code
COPY . .

# Build the Go application
# CGO_ENABLED=0 disables CGO, creating a statically linked binary
# GOOS=linux ensures the binary is built for a Linux environment
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd

# Stage 2: Runner
# Use a minimal base image for the final runtime image
FROM scratch

# Copy the built binary from the builder stage
COPY --from=builder /server /server

# Copy CA certificates and timezone data for HTTPS and time handling
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Expose the port your Go server listens on (e.g., 8080)
EXPOSE 8080

# Command to run the compiled Go application
USER 65532:65532
CMD ["/server"]