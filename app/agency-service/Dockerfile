ARG GOLANG_VERSION

FROM golang:${GOLANG_VERSION:-1.24} AS builder

WORKDIR /app

COPY . .

WORKDIR /app/agency-service

RUN go mod tidy
RUN CGO_ENABLED=0 go build -o agency-service .

FROM alpine:latest
WORKDIR /root/

COPY --from=builder /app/agency-service/agency-service .

CMD ["./agency-service"]