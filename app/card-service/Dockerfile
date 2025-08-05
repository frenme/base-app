ARG GOLANG_VERSION

FROM golang:${GOLANG_VERSION:-1.24} AS builder

WORKDIR /app

COPY ../card-service .

WORKDIR /app/card-service

RUN go mod tidy
RUN CGO_ENABLED=0 go build -o order-service .

FROM alpine:latest
WORKDIR /root/

COPY --from=builder /app/card-service/card-service .

CMD ["./card-service"]