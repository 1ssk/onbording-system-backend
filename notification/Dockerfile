FROM golang:1.24-alpine

RUN apk --no-cache add ca-certificates gcc g++ libc-dev

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /bin/app ./cmd/main.go

EXPOSE 8084 50053

CMD ["/bin/app"]