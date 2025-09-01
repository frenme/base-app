ARG GOLANG_VERSION

FROM golang:${GOLANG_VERSION:-1.24}-alpine AS builder

WORKDIR /app

COPY . .

WORKDIR /app/temp-service

RUN go mod tidy
RUN CGO_ENABLED=0 go build -o temp-service .

FROM alpine:latest
WORKDIR /root/

COPY --from=builder /app/temp-service/temp-service .

CMD ["./temp-service"]