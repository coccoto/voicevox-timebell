# Use Go base image
FROM golang:1.21-alpine

# Install required packages
RUN apk add --no-cache alsa-utils ca-certificates wget

# Set working directory
WORKDIR /app

# Copy backend files
COPY backend/ ./

# Copy frontend files
COPY frontend/ ./frontend/

# Initialize go module if not exists
RUN if [ ! -f go.mod ]; then go mod init voicevox-timebell; fi

# Build the application
RUN go build -o server .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./server"]
