# Stage 1: Build
FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /hrd_backend ./cmd/main.go

# Stage 2: Runtime
FROM alpine:3.18
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app

COPY --from=builder /hrd_backend ./hrd_backend
COPY migrations ./migrations

RUN mkdir -p /app/uploads/questions

EXPOSE 8080
CMD ["./hrd_backend"]
