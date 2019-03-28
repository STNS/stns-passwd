VERSION:=$(shell cat version.go | grep -i version | awk -F= '{print $$2}' | sed -e 's/"//g' | tr -d ' ')
build:
	rm -rf release/*
	gox -os="darwin linux" -arch="386 amd64" -output "release/stns_{{.OS}}_{{.Arch}}/{{.Dir}}"

release:
	git tag -a $(VERSION) -m "bump to $(VERSION)" || true
	goreleaser --rm-dist 
.PHONY: release
