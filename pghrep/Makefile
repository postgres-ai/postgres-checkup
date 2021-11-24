.DEFAULT_GOAL = all

BINARY = pghrep
GOARCH = amd64

VERSION?=0.1
BUILD_TIME?=$(shell date -u '+%Y%m%d-%H%M')
COMMIT?=no #$(shell git rev-parse HEAD)
BRANCH?=no #$(shell git rev-parse --abbrev-ref HEAD)

# Symlink into GOPATH
GITHUB_USERNAME=dmius
BUILD_DIR=${GOPATH}/src/github.com/${GITHUB_USERNAME}/${BINARY}
#BUILD_DIR=${GOPATH}/${BINARY}

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags "-s -w \
	-X version.version=${VERSION} \
	-X version.commit=${COMMIT} \
	-X version.branch=${BRANCH}\
	-X version.buildTime=${BUILD_TIME}"

PLUGINS_SRC := $(wildcard plugins/*.go)

# Build the project
all: clean vet main

main:
	GOARCH=${GOARCH} go build ${LDFLAGS} -o bin/${BINARY} ./cmd/checkup

test:
	go test ./internal/***

vet:
	go vet ./internal/...

fmt:
	go fmt $$(go list ./... | grep -v /vendor/)

clean:
	-rm -f bin/*

run:
	go run ${LDFLAGS} ./cmd/checkup/main.go

.PHONY: all main test vet fmt clean run

