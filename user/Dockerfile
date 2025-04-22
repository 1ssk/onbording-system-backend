FROM golang:1.23.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -o main ./main.go

FROM alpine:latest  

WORKDIR /app

COPY --from=builder /app/main /app/main


COPY .env .env

EXPOSE 8080

CMD ["/app/main"]