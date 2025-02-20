FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /app/shared-service main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/shared-service .

RUN chmod +x /root/shared-service

EXPOSE 8083

CMD ["./shared-service"]