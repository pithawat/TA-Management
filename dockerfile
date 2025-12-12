# ----------------------------------------
# STAGE 1: Base Builder (For compiling the application)
# ----------------------------------------
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Install git for module dependency fetching (if needed)
RUN apk add --no-cache git 

# Copy module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build the final binary
COPY . .
# Compile the application for production
RUN CGO_ENABLED=0 go build -o ta-management .


# ----------------------------------------
# STAGE 2: Test Builder (Includes testing dependencies)
# ----------------------------------------
FROM golang:1.24-alpine AS test-builder
WORKDIR /app


# Copy module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy all application code for running tests
COPY . .

# Install PostgreSQL client for healthchecks/debugging within the test container (optional)
# RUN apk add --no-cache postgresql-client

# ----------------------------------------
# STAGE 3: Final Production Image (Minimal, Secure Runtime)
# ----------------------------------------
FROM alpine:latest AS final
WORKDIR /app

# Security: Ensure only the built binary and necessary files are copied
# Copy the compiled binary from the builder stage
COPY --from=builder /app/ta-management .

# Expose the application port (adjust if necessary)
EXPOSE 8080

# Run the final application
CMD ["./ta-management"]