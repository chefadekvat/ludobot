FROM golang:1.25.3-alpine AS builder

RUN apk add --no-cache git build-base

WORKDIR /bot
COPY . .

RUN make bot

FROM alpine:latest

WORKDIR /root

COPY --from=builder /bot/build/ludobot .
COPY --from=builder /bot/config ./config

CMD ["./ludobot", "--localization-config", "./config/localization/config.yaml", "--bot-config", "./config/bot/config.yaml"]