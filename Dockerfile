FROM golang:1.14-alpine as build
COPY ./pghrep /go/pghrep
RUN apk add --update --no-cache make git
RUN cd /go/pghrep && make install main

FROM alpine:3.11 as production
RUN apk add --update --no-cache \
  bash \
  openssh-client \
  postgresql-client \
  jq \
  curl \
  gawk \
  sed
WORKDIR ./checkup
COPY --from=build /go/pghrep/bin/pghrep ./pghrep/bin/
COPY . .
