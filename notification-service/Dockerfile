FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum .env ./
RUN go mod download

COPY . .
RUN go build -o notification-service ./cmd

FROM alpine:3.22.1
WORKDIR /root/
COPY --from=builder /app/notification-service /app/.env ./

EXPOSE 8005
CMD ["./notification-service"]
