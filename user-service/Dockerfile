FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum .env ./
RUN go mod download

COPY . .
RUN go build -o user-service ./cmd

FROM alpine:3.22.1
WORKDIR /root/
COPY --from=builder /app/user-service /app/.env ./

EXPOSE 8001
CMD ["./user-service"]
