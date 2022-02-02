FROM golang:1.17.1 as builder

ADD blockchain_ws /blockchain_ws

WORKDIR /blockchain_ws

RUN go build -o ws


FROM alpine:3.14

RUN apk add libc6-compat
RUN mkdir /app
WORKDIR /app

COPY --from=builder /blockchain_ws/ws .

EXPOSE 5000

CMD ./ws