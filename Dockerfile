# Dockerfile
FROM golang:1.25-alpine

WORKDIR /app

ENV GOTOOLCHAIN=auto

COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build
RUN go build -o task-scheduler .

EXPOSE 8080
CMD ["./task-scheduler"]
