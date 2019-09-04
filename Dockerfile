FROM alpine:3.9

ARG CHECKUP_HOSTS
ARG CHECKUP_SNAPSHOT_DISTANCE
ARG CHECKUP_CONFIG_PATH

ENV CHECKUP_HOSTS=$CHECKUP_HOSTS
ENV CHECKUP_SNAPSHOT_DISTANCE=$CHECKUP_SNAPSHOT_DISTANCE
ENV CHECKUP_CONFIG_PATH=$CHECKUP_CONFIG_PATH

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
