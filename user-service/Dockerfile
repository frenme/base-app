FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /app/user-service main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/user-service .

RUN chmod +x /root/user-service

EXPOSE 8081

CMD ["./user-service"]