FROM golang:1.17.7-alpine3.15 AS build-env

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build -o ./secret-scanner

RUN apk add --no-cache git postgresql-client

EXPOSE 8080

CMD ["./secret-scanner"]
