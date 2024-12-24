FROM golang:1.19-alpine

WORKDIR /app

# Install required system dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod and sum files
COPY go.mod ./

# Initialize go.mod and download dependencies
RUN go mod download
RUN go mod tidy

# Copy the source code
COPY . .

# Run go mod tidy again after copying source code
RUN go mod tidy

# Build the application
RUN go build -o main .

# Create data directory for SQLite
RUN mkdir -p /app/data

# Expose port 8080
EXPOSE 8080

# Command to run the application
CMD ["./main"] 