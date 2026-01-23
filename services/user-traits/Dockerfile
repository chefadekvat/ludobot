FROM golang:1.25.3-alpine AS builder

ARG USER_TRAITS_ROOT
ARG SQLGEN_ROOT

RUN apk add --no-cache git build-base

COPY ${USER_TRAITS_ROOT} /app
COPY ${SQLGEN_ROOT}/main.go /sqlgen/main.go

ENV SQLGEN_ROOT=/sqlgen

WORKDIR /app
RUN make user-traits

FROM alpine:latest

WORKDIR /root

COPY --from=builder /app/build/user-traits .