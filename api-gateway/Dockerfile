FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o api-gateway ./cmd

FROM alpine:3.22.1
WORKDIR /root/
COPY --from=builder /app/api-gateway .

EXPOSE 8080
CMD ["./api-gateway"]
