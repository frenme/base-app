ARG GOLANG_VERSION

FROM golang:${GOLANG_VERSION:-1.24} AS builder

WORKDIR /app

COPY . .

WORKDIR /app/order-service

RUN go mod tidy
RUN CGO_ENABLED=0 go build -o order-service .

FROM alpine:latest
WORKDIR /root/

COPY --from=builder /app/order-service/order-service .

CMD ["./order-service"]