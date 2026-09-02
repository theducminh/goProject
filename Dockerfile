# Multi-stage build
# Stage 1: Build
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Cài đặt dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go.mod và go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Generate Swagger docs
ENV PATH="/root/go/bin:${PATH}"
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.6 && \
    swag init -g main.go --parseDependency --output docs

# Build binary
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main .

# Stage 2: Runtime
FROM alpine:latest

WORKDIR /app

# Cài đặt ca-certificates và timezone data
RUN apk --no-cache add ca-certificates tzdata

RUN addgroup -S app && adduser -S -G app app

# Copy binary từ builder
COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs

USER app

# Expose port
EXPOSE 8080
EXPOSE 9090

# Run
CMD ["./main"]
