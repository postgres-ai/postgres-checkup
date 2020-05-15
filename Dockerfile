FROM golang:1.14 as build
COPY ./pghrep /go/pghrep
RUN cd /go/pghrep && make install main

FROM alpine:3.9 as production
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
