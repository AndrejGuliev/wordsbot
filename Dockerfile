FROM golang:1.18.4-alpine3.16 AS builder

RUN go version

COPY . /github.com/AndrejGuliev/wordsbot
WORKDIR /github.com/AndrejGuliev/wordsbot

RUN go mod download
RUN GOOS=linux go build -o ./.bin/bot ./main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 github.com/AndrejGuliev/wordsbot/.bin/bot .
COPY --from=0 github.com/AndrejGuliev/wordsbot/configs configs/

EXPOSE 3306

CMD ["./bot"]