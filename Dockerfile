# Stage 1: build binary
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary statically linked for Linux
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server .

# Stage 2: runtime
FROM alpine:3.20

# Install ca-certificates (for HTTPS/Mailtrap) and tzdata (for Asia/Jakarta timezone)
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server /app/server

# Copy badword dictionary file needed by config.LoadBadWords
COPY --from=builder /app/config/badword.txt /app/config/badword.txt

EXPOSE 8080

CMD ["/app/server"]
