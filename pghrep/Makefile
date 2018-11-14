.DEFAULT_GOAL = all

BINARY = pghrep
VET_REPORT = vet.report
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

# Build the project
all: clean vet linux freebsd darwin windows

linux:
	GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o bin/${BINARY} ./src/

freebsd:
#	GOOS=freebsd GOARCH=${GOARCH} go build ${LDFLAGS} -o bin/${BINARY}-freebsd-${GOARCH} ./src/

darwin:
#	GOOS=darwin GOARCH=${GOARCH} go build ${LDFLAGS} -o bin/${BINARY}-darwin-${GOARCH} ./src/

windows:
#	GOOS=windows GOARCH=${GOARCH} go build ${LDFLAGS} -o bin/${BINARY}.exe .

vet:
	go vet ./... > ${VET_REPORT} 2>&1

fmt:
	go fmt $$(go list ./... | grep -v /vendor/)

clean:
	-rm -f ${TEST_REPORT}
	-rm -f ${VET_REPORT}
	-rm -f bin/${BINARY}*

run:
	go run ${LDFLAGS} ./src/*

.PHONY: all linux darwin windows vet fmt clean run
