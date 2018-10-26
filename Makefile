.PHONY: vendor_install vendor_status vendor_update vendor_sync install build doc fmt gen lint test vet godep bench

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

vendor_install:
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

vendor_status:
	dep status

vendor_update:
	dep ensure -update

vendor_sync:
	dep ensure

install: vendor_sync
	go get -x $(GOBUILD_VERSION_ARGS) -t -v ./...

build: vendor_sync
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
# curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $GOPATH/bin v1.10.2
lint:
	golangci-lint run

test:
	go test ./...

# Runs benchmarks
bench:
	go test ./... -bench=.

# https://godoc.org/golang.org/x/tools/cmd/vet
vet:
	go vet ./...
