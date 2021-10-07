FROM golang:1.17.1 as builder

ADD tsundoku_blockchain_ws /tsundoku_blockchain_ws

WORKDIR /tsundoku_blockchain_ws

RUN go build -o ws


FROM alpine:3.14

RUN apk add libc6-compat
RUN mkdir /app
WORKDIR /app

COPY --from=builder /tsundoku_blockchain_ws/ws .

EXPOSE 5000

CMD ./ws