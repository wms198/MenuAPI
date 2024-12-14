FROM golang:1.23.4-alpine3.20 AS builder

RUN mkdir /app
WORKDIR /app

RUN apk add --update-cache build-base

COPY . .
RUN go mod download

ENV CGO_ENABLED=1
RUN go build -o /usr/local/bin/goapi .

FROM alpine
COPY .env .
COPY --from=builder /usr/local/bin/goapi /usr/local/bin/goapi

CMD ["goapi"]