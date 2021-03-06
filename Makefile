.PHONY: install build doc fmt gen lint test vet bench

PKG_NAME=$(shell basename `pwd`)
TARGET_OS="linux"
VERSION_VAR=main.BuildVersion
TIMESTAMP_VAR=main.BuildTimestamp
VERSION=$(shell git describe --always --dirty --tags)
TIMESTAMP=$(shell date -u '+%Y-%m-%d_%I:%M:%S%p')
GOBUILD_VERSION_ARGS := -ldflags "-X $(VERSION_VAR)=$(VERSION) -X $(TIMESTAMP_VAR)=$(TIMESTAMP)"
CONTAINER_NAME=$(shell grep "FROM " contrib/docker/Dockerfile | sed 's/FROM \(.*\).*\:.*/\1/')
CONTAINER_VERSION=$(shell grep "FROM " contrib/docker/Dockerfile | sed 's/FROM .*\:\(.*\).*/\1/')

default: install

install:
	go get -x $(GOBUILD_VERSION_ARGS) -t -v ./...

build:
	go build -x $(GOBUILD_VERSION_ARGS) -v -o ./bin/$(PKG_NAME)

docker: docker_build docker_container

docker_container:
	docker build -t $(CONTAINER_NAME)-kube-template:$(CONTAINER_VERSION) ./contrib/docker

docker_build:
	CGO_ENABLED=0 GOOS=$(TARGET_OS) go build -x $(GOBUILD_VERSION_ARGS) -a -installsuffix cgo -v -o ./contrib/docker/bin/$(PKG_NAME)

clean:
	rm -dRf ./bin

doc:
	godoc -http=:6060

fmt:
	go fmt ./...

gen:
	go generate

# https://golangci.com/
# curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.27.0
lint:
	golangci-lint run --timeout=300s

test:
	go test ./...

# Runs benchmarks
bench:
	go test ./... -bench=.

# https://godoc.org/golang.org/x/tools/cmd/vet
vet:
	go vet ./...
