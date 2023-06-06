FROM golang:1.20.4-alpine AS builder

WORKDIR /app
COPY . .
RUN apk add git --no-cache && go build -o ./proxy

FROM alpine:3.18.0

RUN apk add ca-certificates --no-cache
WORKDIR /app
COPY --from=builder /app/proxy /app/proxy

ENTRYPOINT [ "/app/proxy" ]
