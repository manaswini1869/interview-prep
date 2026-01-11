FROM golang:1.23

# Copy everything first (Bad practice: breaks cache on any code change)
COPY . /app
WORKDIR /app

# Download deps every time?
RUN go mod download

# Build the binary
RUN go build -o worker-registry main.go

# Run as root? (Security risk)
CMD ["./worker-registry"]