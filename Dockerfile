FROM alpine:3.9

ARG CHECKUP_HOSTS
ARG CHECKUP_SNAPSHOT_DISTANCE
ARG CHECKUP_CONFIG_PATH

ADD run_checkup.sh

RUN apk add --update --no-cache \
    bash \
    openssh-client \
    postgresql-client \
    coreutils \
    jq \
    curl \
    go \
    gawk \
    sed \
    make \
    build-base \
    git

COPY . .
