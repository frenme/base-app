FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 go build -o /app/order-service main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/order-service .

CMD ["./order-service"]