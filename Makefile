SHELL := /bin/bash

ifndef RHM_VERSION
RHM_VERSION := $(shell git describe --tags --abbrev=14)
endif
LDFLAGS := -X main.Version=$(RHM_VERSION)

.PHONY: all
all:
	@go install -v -ldflags '$(LDFLAGS)'

.PHONY: clean
clean:
	@-go clean -i

.PHONY: ci
ci: check-gofmt check-goimports check-golint vet test test-race

# goimports doesn't support the -s flag to simplify code, therefore we use both
# goimports and gofmt -s.
.PHONY: check-gofmt
check-gofmt:
	diff <(gofmt -s -d .) <(printf "")

.PHONY: check-goimports
check-goimports:
	diff <(goimports -d .) <(printf "")

.PHONY: check-golint
check-golint:
	diff <(golint ./... | grep -v vendor/) <(printf "")

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test:
	go test -v -cpu=2 `go list ./... | grep -v /vendor/`

.PHONY: test-race
test-race:
    go test -v -cpu=1,2,4 -short -race `go list ./... | grep -v /vendor/`

.PHONY: test-with-vendor 
test-with-vendor:
	go test -v -cpu=2 ./...