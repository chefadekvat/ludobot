FROM alpine:latest AS builder

ARG USER_TRAITS_ROOT
ARG SQLGEN_ROOT

RUN apk add --no-cache git build-base go

COPY ${USER_TRAITS_ROOT} /app
COPY ${SQLGEN_ROOT}/main.go /sqlgen/main.go

ENV SQLGEN_ROOT=/sqlgen

WORKDIR /app
RUN make user-traits

FROM alpine:latest

RUN apk add --no-cache python3 py3-virtualenv \
    postgresql18 \
    postgresql18-client \
    libpq-dev \
    python3-dev \
    go

RUN mkdir -p /opt/yandex
RUN python3 -m virtualenv --python=/usr/bin/python3 /opt/yandex/testsuite

RUN /opt/yandex/testsuite/bin/pip install yandex-taxi-testsuite[postgresql,postgresql-binary]
RUN /opt/yandex/testsuite/bin/pip install asyncpg

WORKDIR /app

COPY --from=builder /app/build/user-traits .
