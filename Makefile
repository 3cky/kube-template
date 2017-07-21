.PHONY: vendor_clean vendor_fetch vendor_update vendor_sync install build doc fmt lint test vet godep bench

PKG_NAME=$(shell basename `pwd`)
TARGET_OS="linux"
CONTAINER_NAME=$(shell grep "FROM " contrib/docker/Dockerfile | sed 's/FROM \(.*\).*\:.*/\1/')
CONTAINER_VERSION=$(shell grep "FROM " contrib/docker/Dockerfile | sed 's/FROM .*\:\(.*\).*/\1/')

default: install

vendor_clean:
	govendor remove +u

vendor_fetch:
	govendor fetch +out

vendor_update:
	govendor update +vendor

vendor_sync:
	govendor sync -v

install: vendor_sync
	go get -t -v ./...

build: vendor_sync
	go build -v -o ./bin/$(PKG_NAME)

docker: docker_build docker_container

docker_container:
	docker build -t $(CONTAINER_NAME)-kube-template:$(CONTAINER_VERSION) ./contrib/docker

docker_build:
	CGO_ENABLED=0 GOOS=$(TARGET_OS) go build -a -installsuffix cgo -v -o ./contrib/docker/bin/$(PKG_NAME)

clean: vendor_clean
	rm -dRf ./bin

doc:
	godoc -http=:6060

fmt:
	go fmt ./...

# https://github.com/golang/lint
# go get github.com/golang/lint/golint
lint:
	golint ./...

test:
	go test ./...

# Runs benchmarks
bench:
	go test ./... -bench=.

# https://godoc.org/golang.org/x/tools/cmd/vet
vet:
	go vet ./...
