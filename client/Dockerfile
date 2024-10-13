# Dockerfile for client
FROM golang:1.20-alpine

WORKDIR /app

COPY . .

# Build the client binary
RUN go mod init client
RUN go mod tidy
RUN go build -o client client.go

# Run the client binary
CMD ["./client"]
