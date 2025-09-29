# Dockerfile
FROM golang:1.25-alpine

WORKDIR /app

# (optional) let Go fetch the right minor toolchain automatically
ENV GOTOOLCHAIN=auto

# Cache deps
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build
RUN go build -o task-scheduler .

EXPOSE 8080
CMD ["./task-scheduler"]
