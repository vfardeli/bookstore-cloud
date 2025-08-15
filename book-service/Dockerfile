FROM golang:1.25-alpine

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o book-service ./cmd

EXPOSE 8000
CMD ["./book-service"]
