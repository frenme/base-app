ARG GOLANG_VERSION

FROM golang:${GOLANG_VERSION:-1.24}-alpine AS builder

WORKDIR /app

COPY . .

WORKDIR /app/user-service

RUN go mod tidy
RUN CGO_ENABLED=0 go build -o user-service .

FROM alpine:latest
WORKDIR /root/

COPY --from=builder /app/user-service/user-service .

CMD ["./user-service"]