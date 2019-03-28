VERSION:=$(shell cat version.go | grep -i version | awk -F= '{print $$2}' | sed -e 's/"//g' | tr -d ' ')
TEST ?= $(shell go list ./... | grep -v -e vendor -e keys -e tmp)
INFO_COLOR=\033[1;34m
RESET=\033[0m
BOLD=\033[1m
ifeq ("$(shell uname)","Darwin")
GO ?= GO111MODULE=on go
else
GO ?= GO111MODULE=on /usr/local/go/bin/go
endif

build:
	rm -rf release/*
	gox -os="darwin linux" -arch="386 amd64" -output "release/stns_{{.OS}}_{{.Arch}}/{{.Dir}}"

release:
	git tag -a $(VERSION) -m "bump to $(VERSION)" || true
	goreleaser --rm-dist

test: ## Run test
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Testing$(RESET)"
	$(GO) test -v $(TEST) -timeout=30s -parallel=4
	$(GO) test -race $(TEST)

.PHONY: release
