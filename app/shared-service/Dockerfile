ARG GOLANG_VERSION

FROM golang:${GOLANG_VERSION:-1.24} AS builder

WORKDIR /app

COPY . .

WORKDIR /app/shared-service

RUN go mod tidy
RUN CGO_ENABLED=0 go build -o shared-service .

FROM alpine:latest
WORKDIR /root/

COPY --from=builder /app/shared-service/shared-service .

CMD ["./shared-service"]