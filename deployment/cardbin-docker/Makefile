GOPATH:=$(shell go env GOPATH)

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux go build -o cardbin ../../cmd/cardbin/main.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t cardbin/cardbin:latest
