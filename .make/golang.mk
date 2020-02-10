# This Makefile describes the behavior of a node that was written on golang.


# Get the go bin path.
GO := $(shell which go)
ifeq ($(GO),)
	# Case when bainary is not installed. We have target below for installing stable version.
	GO := /usr/local/go/bin/go
	GO_VERSION := 1.13
	GO_RELEASE_LINK := https://dl.google.com/go/go$(GO_VERSION).linux-amd64.tar.gz
	# Set golang global env variables.
	export GOROOT=/usr/local/go
	export GOPATH=/go
endif


# Tell git to use ssh clone instead https (for private repos).
# Workaround with invoking at initialization stage of this Makefile.
GITHUB_REWRITE := $(shell git config --global url.git@github.com:.insteadOf https://github.com/)


# Find source files (using find instead of wildcard resolves any depth issue).
GO_SOURCE_FILES := $(shell find * -type f -name '*.go')
GO_TEST_SOURCE_FILES := $(shell find * -type f -name '*_test.go')


$(GO):
	#### Node( '$(NODE)' ).Call( '$@' )
	@mkdir /lib64
	@ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
	wget $(GO_RELEASE_LINK) -qO- | tar -C /usr/local -xzf -
	$(GO) version


.PHONY: golang-fmt
golang-fmt: $(GO) $(GO_SOURCE_FILES)
	#### Node( '$(NODE)' ).Call( '$@' )
	$(GO) fix ./...
	$(GO) fmt ./...
	$(GO) vet ./...


# Performance section. Testing, benchmarking, profiling, tracing, debugging, etc.
.PHONY: golang-test
golang-test: $(GO) $(GO_TEST_SOURCE_FILES)
	#### Node( '$(NODE)' ).Call( '$@' )
	$(GO) test -race -cover ./...


# Build and run section. Convert source code to executable and provide process.
# Provide multiple options for building (bin, lib, etc).

# Calculate build variables.
TIMESTAMP := $(shell date +%s)
VERSION := LOCAL
# TODO: BELOW
# define LDFLAGS
# -w \
# -linkmode external \
# -extldflags '-static' \
# -X 'main.NODE=$(BASE)' \
# -X 'main.VERSION=$(VERSION)' \
# -X 'main.TIMESTAMP=$(TIMESTAMP)'
# endef
# # CGO_ENABLED=1 \
# # CC='gcc' \
# -ldflags "$(LDFLAGS)" \
# TODO: ABOVE


# Targets for building executable binaries.
$(BASE)-Linux-x86_64: $(GO) $(GO_SOURCE_FILES)
	#### Node( '$(NODE)' ).Call( '$@' )
	GOOS=linux GOARCH=amd64 $(GO) build \
		-v \
		-o $(BASE)-Linux-x86_64


$(BASE)-Darwin-x86_64: $(GO) $(GO_SOURCE_FILES)
	#### Node( '$(NODE)' ).Call( '$@' )
	GOOS=darwin GOARCH=amd64 $(GO) build \
		-v \
		-o $(BASE)-Darwin-x86_64


# Targets for building golang plugin libraries.
$(BASE)-Linux-x86_64-plugin.so: $(GO) $(GO_SOURCE_FILES)
	#### Node( '$(NODE)' ).Call( '$@' )
	GOOS=linux GOARCH=amd64 $(GO) build \
		-v \
		-buildmode=plugin \
		-o $(BASE)-Linux-x86_64-plugin.so

##############################
### GOLANG MODULES SECTION ###
##############################


go.mod: $(GO)
	#### Node( '$(NODE)' ).Call( '$@' )
	$(GO) mod init


.PHONY: golang-mod-tidy
golang-mod-tidy: $(GO) go.mod
	#### Node( '$(NODE)' ).Call( '$@' )
	$(GO) mod tidy


.PHONY: golang-mod-graph
golang-mod-graph: $(GO) go.mod
	#### Node( '$(NODE)' ).Call( '$@' )
	$(GO) mod graph


# TODO:
# https://blog.golang.org/using-go-modules

# This post introduced these workflows using Go modules:
#
# go mod init creates a new module, initializing the go.mod file that describes it.
# go build, go test, and other package-building commands add new dependencies to go.mod as needed.
# go list -m all prints the current module’s dependencies.
# go get changes the required version of a dependency (or adds a new dependency).
# go mod tidy removes unused dependencies.

# TODO: deal with
# Imports and canonical module paths
