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
	-X main.version=${VERSION} \
	-X main.commit=${COMMIT} \
	-X main.branch=${BRANCH}\
	-X main.buildTime=${BUILD_TIME}"

PLUGINS_SRC := $(wildcard plugins/*.go)

# Build the project
all: clean install vet buildplugins main

install:
	go get golang.org/x/net/html
	go get github.com/dustin/go-humanize

buildplugins:
	@for f in $(basename $(patsubst plugins/%.go,%,$(PLUGINS_SRC))); do \
		go build ${LDFLAGS} -buildmode=plugin -o bin/$$f.so ./plugins/$$f.go ; \
	done

main:
	GOARCH=${GOARCH} go build ${LDFLAGS} -o bin/${BINARY} ./src/

test:
	go test ./src/
	go test ./src/fmtutils/
	go test ./src/checkup/*/

vet:
	# Command go vet examine all go files in a single scope, which isn't suitable for plugins.
	# We will run it separately for plugins and main sources.
	@for f in $(basename $(patsubst plugins/%.go,%,$(PLUGINS_SRC))); do \
		go vet ./plugins/$$f.go || exit 1 ; \
	done
	go vet ./src/...

fmt:
	go fmt $$(go list ./... | grep -v /vendor/)

clean:
	-rm -f bin/*

run:
	go run ${LDFLAGS} ./src/*

.PHONY: all buildplugins main test vet fmt clean run

