FROM golang:1.17.7-alpine3.15 AS build-env

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build -o ./secret-scanner

## ADD git to the image (for repos downloading)
RUN apk add --no-cache git

## Add the wait script to the image
ENV WAIT_VERSION 2.9.0
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/$WAIT_VERSION/wait /wait
RUN chmod +x /wait

EXPOSE 8080

CMD /wait && ./secret-scanner
