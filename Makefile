# Copyright 2018 The ZJU-SEL Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

GO ?= go
PROJECT := github.com/ZJU-SEL/capstan
BINDIR := /usr/local/bin
ifeq ($(GOPATH),)
export GOPATH := $(CURDIR)/_output
unexport GOBIN
endif
GOBINDIR := $(word 1,$(subst :, ,$(GOPATH)))
PATH := $(GOBINDIR)/bin:$(PATH)
GOPKGDIR := $(GOPATH)/src/$(PROJECT)
GOPKGBASEDIR := $(shell dirname "$(GOPKGDIR)")


.PHONY: all
all: build

.PHONY: help
help:
	@echo "Usage: make <target>"
	@echo
	@echo " * 'install' - Install capstan to system locations."
	@echo " * 'uninstall' - Uninstall capstan to system locations."
	@echo " * 'build' - Build capstan."
	@echo " * 'clean' - Clean artifacts."

.PHONY: check-gopath
check-gopath:
ifeq ("$(wildcard $(GOPKGDIR))","")
	mkdir -p "$(GOPKGBASEDIR)"
	ln -s "$(CURDIR)" "$(GOPKGBASEDIR)/capstan"
endif
ifndef GOPATH
	$(error GOPATH is not set)
endif

.PHONY: build
build: check-gopath
		$(GO) install \
		$(PROJECT)/cmd/capstan

.PHONY: clean
clean:
	find . -name \*~ -delete
	find . -name \#\* -delete

.PHONY: install
install: check-gopath
	install -D -m 755 $(GOBINDIR)/bin/capstan $(BINDIR)/capstan

.PHONY: uninstall
uninstall:
	rm -f $(BINDIR)/capstan

.PHONY: fmt
fmt: check-gopath
	files=$$(cd $(GOPKGDIR) && find . -not \(  \( -wholename '*/vendor/*' \) -prune \) -name '*.go' | xargs gofmt -s -l | tee >(cat - >&2)); [ -z "$$files" ]

.PHONY: lint
lint:
	hack/verify-gofmt.sh
	hack/verify-govet.sh
	hack/verify-boilerplate.sh

.PHONY: install.tools
install.tools:
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install

.PHONY: \
	help \
	check-gopath \
	build \
	clean \
	install \
	uninstall \
	lint \
	fmt \
	install.tools