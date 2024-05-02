FROM golang:1.17-alpine as build
COPY ./pghrep /go/pghrep
RUN apk add --update --no-cache make
RUN cd /go/pghrep && make main

FROM alpine:3.11 as production
RUN apk add --update --no-cache \
  bash \
  openssh-client \
  postgresql-client \
  jq \
  curl \
  gawk \
  sed
RUN wget https://github.com/jqlang/jq/releases/download/jq-1.7.1/jq-linux-amd64
RUN cp ./jq-linux-amd64 /usr/local/bin/jq && chmod +x /usr/local/bin/jq 
WORKDIR ./checkup
COPY --from=build /go/pghrep/bin/pghrep ./pghrep/bin/
COPY . .
