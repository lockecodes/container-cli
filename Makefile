SHELL := bash
PODMAN ?= false
DOCKER:="$(shell if ${PODMAN}; then \
		echo podman; \
		else echo docker; \
	fi \
)"
GO_RELEASER_VERSION := v2.6.1

format:
	go fmt ./...

test-release:
	$(DOCKER) run \
	  --rm -v $(shell pwd):/go/src/gitlab.com/locke-codes/container-cli \
	  -w /go/src/gitlab.com/locke-codes/container-cli \
	  goreleaser/goreleaser:$(GO_RELEASER_VERSION) \
	  release \
		--clean \
		--auto-snapshot \
		--skip publish
