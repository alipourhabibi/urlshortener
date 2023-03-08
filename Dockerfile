FROM golang:1.20.1-alpine3.16 AS build

LABEL Project="urlshortener"

ENV GO111MODULE=on

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build  -o shortener
CMD ["./shortener", "run"]
