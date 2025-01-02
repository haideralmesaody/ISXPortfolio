# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install required system dependencies including git
RUN apk add --no-cache gcc musl-dev sqlite-dev git

# Set CGO flags for proper SQLite compilation
ENV CGO_CFLAGS="-D_LARGEFILE64_SOURCE"
ENV CGO_ENABLED=1

# Copy go mod files first for better caching
COPY go.mod go.sum* ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /app

# Install Chrome and dependencies
RUN apk add --no-cache \
    chromium \
    chromium-chromedriver \
    sqlite \
    nss \
    freetype \
    freetype-dev \
    harfbuzz \
    ca-certificates \
    ttf-freefont \
    dbus \
    udev \
    xvfb \
    # Add these for better browser support
    fontconfig \
    xvfb-run

# Set correct Chrome paths and shared memory
ENV CHROME_BIN=/usr/bin/chromium \
    CHROME_PATH=/usr/lib/chromium/ \
    DISPLAY=:99

# Fix shared memory issues
RUN mkdir -p /dev/shm && chmod 1777 /dev/shm

# Ensure data directories exist with proper permissions
RUN mkdir -p /app/data && \
    chmod 777 /app/data

RUN mkdir -p /app/data/pdfs && \
    chmod 777 /app/data/pdfs

# Copy the binary from builder
COPY --from=builder /app/main .

# Start Xvfb and the application
CMD ["sh", "-c", "echo 'Starting application...' && Xvfb :99 -screen 0 1024x768x16 & sleep 1 && ./main"] 