EXECUTABLES := app

GO ?= go
GOFMT ?= gofmt "-s"
GOIMPORTS ?= goimports "-s"

.PHONY: lint
lint:
	golangci-lint run
