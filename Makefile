.PHONY: build doc fmt lint test vet godep install bench

PKG_NAME=$(shell basename `pwd`)

default: build

install:
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
