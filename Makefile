checks: build test pre-commit
.PHONY: checks

build: oe
.PHONY: build

oe:
	go build ./cmd/$@
.PHONY: oe

test:
	go test -v ./... -coverpkg=./... -coverprofile=.coverprofile
.PHONY: test

clean:
	rm -f oe .coverprofile
.PHONY: clean

pre-commit:
	pre-commit run -a
.PHONY: pre-commit

gocovsh:
	go test -v ./... -coverpkg=./... -coverprofile=.coverprofile
	gocovsh --profile .coverprofile
	rm .coverprofile
.PHONY: gocovsh

ci-tools:
	go install honnef.co/go/tools/cmd/staticcheck@c88a4f6cc653658bc2ac942ae7a450a7c14d08bd # v0.7.0 (2026.1)
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@c0d3ddc9cf3faa61a4e378e879ece580256d76e5 # v2.12.2
	go install github.com/securego/gosec/v2/cmd/gosec@5e5517beec77b8228ba43ec8d7cc22d82ed31924 # v2.25.0
.PHONY: ci-tools
