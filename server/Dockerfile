# Dockerfile for server
FROM golang:1.20-alpine

# fix error Failed to initialize database: failed to create table: Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
# https://kodekloud.com/community/t/issue-while-running-go-webapp-sample-application-in-local-machine/363181

WORKDIR /app

COPY . .

# Install SQLite
RUN apk add --no-cache sqlite


# Build the server binary
RUN go mod init server
RUN go mod tidy
RUN go build -o server server.go

# Expose port 8080
EXPOSE 8080

# Run the server binary
CMD ["./server"]
