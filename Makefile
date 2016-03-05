.PHONY: vendor_clean vendor_restore vendor_get install build doc fmt lint test vet godep bench

PKG_NAME=$(shell basename `pwd`)

GOPATH := ${PWD}/_vendor:${GOPATH}
export GOPATH

default: build

vendor_clean:
	rm -dRf ./_vendor/src

vendor_fetch:
	GOPATH=${PWD}/_vendor go get -d -u -v \
	    github.com/golang/glog \
	    github.com/spf13/cobra \
	    github.com/spf13/viper \
	    k8s.io/kubernetes/pkg/api \
	    k8s.io/kubernetes/pkg/labels \
	    k8s.io/kubernetes/pkg/util \
	    k8s.io/kubernetes/pkg/client/unversioned

vendor_restore:
	cd ${PWD}/_vendor/src/k8s.io/kubernetes && \
	    GOPATH=${PWD}/_vendor godep restore
#	mkdir -p _vendor/src/github.com/3cky
#	ln -s `pwd` _vendor/src/github.com/3cky/kube-template # for IDEA Go plugin project build and run

vendor_get: vendor_clean vendor_fetch vendor_restore

install: vendor_get
	go get -t -v ./...

build: vet
	go build -v -o ./bin/$(PKG_NAME)

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

godep:
	godep save ./...
