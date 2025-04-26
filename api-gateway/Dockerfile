FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -o main ./main.go

FROM alpine:latest  

WORKDIR /app

COPY --from=builder /app/main /app/main


COPY .env .env

EXPOSE 8000

CMD ["/app/main"]