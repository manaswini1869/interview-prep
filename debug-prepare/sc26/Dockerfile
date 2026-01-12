FROM golang:latest
WORKDIR /app
COPY . .
RUN go build -o server main.go
CMD ["./server"]