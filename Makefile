SHELL := /bin/bash
VERSION := v0.0.1
NAME := rhm
#ifndef RHM_VERSION
#RHM_VERSION := $(shell git describe --tags --abbrev=14)
#endif
#LDFLAGS := -X main.Version=$(RHM_VERSION)

.PHONY: all
all:
	@go install -v

.PHONY: clean
clean:
	@-go clean -i

.PHONY: ci
ci: check-gofmt check-golint vet test test-race

# goimports doesn't support the -s flag to simplify code, therefore we use both
# goimports and gofmt -s.
.PHONY: check-gofmt
check-gofmt:
	diff <(gofmt -d -s .) <(printf "")

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

build:
	go build -ldflags "-X main.Version=$(VERSION)"

deps:
	go get github.com/c4milo/github-release
	go get github.com/mitchellh/gox

compile:
	@rm -rf build/
	@gox -ldflags "-X main.Version=$(VERSION)" \
	-osarch="darwin/amd64" \
	-osarch="linux/amd64" \
	-os="windows" \
	-output "build/{{.Dir}}_$(VERSION)_{{.OS}}_{{.Arch}}/$(NAME)" \
	./...

dist: compile
	$(eval FILES := $(shell ls build))
	@rm -rf dist && mkdir dist
	@for f in $(FILES); do \
		(cd $(shell pwd)/build/$$f && tar -cvzf ../../dist/$$f.tar.gz *); \
		(cd $(shell pwd)/dist && shasum -a 512 $$f.tar.gz > $$f.sha512); \
		echo $$f; \
	done

release: dist
	@latest_tag=$$(git describe --tags `git rev-list --tags --max-count=1`); \
	comparison="$$latest_tag..HEAD"; \
	if [ -z "$$latest_tag" ]; then comparison=""; fi; \
	changelog=$$(git log $$comparison --oneline --no-merges --reverse); \
	github-release feedhenry/$(NAME) $(VERSION) "$$(git rev-parse --abbrev-ref HEAD)" "**Changelog**<br/>$$changelog" 'dist/*';