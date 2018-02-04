# Makefile adapted from the awesome one found @ apex/up :)

GO ?= go

vendor:
	@dep ensure
.PHONY: vendor

# Build all files.
build:
	@echo "==> Building"
	@$(GO) generate ./...
.PHONY: build

# Install from source.
install:
	@echo "==> Installing tldrfeed ${GOPATH}/bin/tldrfeed"
	@$(GO) install ./...
.PHONY: install

testsetup:
	-@docker rm -f tldrfeed-test &> /dev/null
	@docker run -d --name tldrfeed-test -p 17017:27017 mvertes/alpine-mongo > /dev/null
	export TLDRFEED_TEST_DB="0.0.0.0:17017"
	export TLDRFEED_TEST_DB_ENABLED=1
.PHONY: testclean

# Run all tests.
test: testsetup
	@$(GO) test ./... && echo "\n==>\033[32m Ok\033[m\n"
.PHONY: test

lint:
	@gometalinter --vendor --exclude ineffassign --exclude errcheck --exclude megacheck ./...
.PHONY: lint

docker:
	@env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -a -installsuffix cgo ./cmd/tldrfeed
	@docker build -t tldrfeed .
.PHONY: docker

# Show source statistics.
cloc:
	@cloc -exclude-dir=vendor .
.PHONY: cloc

# Show to-do items per file.
todo:
	@rg TODO:
.PHONY: todo

tools:
	brew install dep
	@$(GO) get -u gopkg.in/alecthomas/gometalinter.v2
	@gometalinter --install
