SERVICE ?= $(PWD)
BUILDENV=
BUILDENV+=CGO_ENABLED=0 
GIT_HASH := $(CIRCLE_SHA1)
GIT_HASH ?= $(shell git rev-parse HEAD)
LINKFLAGS :=-s -X main.gitHash=$(GIT_HASH) -extldflags "-static"
TESTFLAGS := -v -cover
LINT_FLAGS :=--disable-all --enable=vet --enable=vetshadow --enable=golint --enable=ineffassign --enable=goconst --enable=gofmt
LINTER := gometalinter.v1

EMPTY :=
SPACE := $(EMPTY) $(EMPTY)
join-with = $(subst $(SPACE),$1,$(strip $2))

LEXC :=
ifdef LINT_EXCLUDE
	LEXC := $(call join-with,|,$(LINT_EXCLUDE))
endif

.DEFAULT_GOAL := rebuild

.PHONY: bootstrap_github_credential
bootstrap_github_credential:
	@echo "machine github.com login $(GH_USERNAME) password $(GH_PASSWORD)" > ~/.netrc

.PHONY: install_packages
install_packages: bootstrap_github_credential
	go get -t -v ./...

.PHONY: install_tools
install_tools:
	go get -u gopkg.in/alecthomas/$(LINTER)
	$(LINTER) --install

.PHONY: lint
lint:
ifdef LEXC
	$(LINTER) --exclude '$(LEXC)' $(LINT_FLAGS) ./...
else
	$(LINTER) $(LINT_FLAGS) ./...
endif

.PHONY: clean
clean:
	rm -f $(SERVICE)

$(SERVICE):
	$(BUILDENV) go build -o $(SERVICE) -a -ldflags '$(LINKFLAGS)' .

.PHONY: test
test:
	$(BUILDENV) go test $(TESTFLAGS) ./...

.PHONY: rebuild
rebuild: clean $(SERVICE);

.PHONY: default
default: rebuild ;

.PHONY: build-container
build-container: bootstrap_github_credential install_packages install_tools lint test rebuild
