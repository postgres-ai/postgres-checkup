FROM golang:1.14 as build
COPY ./pghrep /go/pghrep
RUN cd /go/pghrep && make install main
#RUN go build ./src/main.go

FROM alpine:3.9 as production
RUN apk add --update --no-cache \
  bash \
  openssh-client \
  postgresql-client \
  #coreutils \
  jq \
  curl \
  gawk \
  sed
  #git
WORKDIR ./checkup
COPY --from=build /go/pghrep/bin/pghrep ./pghrep/bin/
COPY . .
