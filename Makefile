# Copyright (c) 2018 The ZJU-SEL Authors.
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
PWD := $(shell pwd)
BASE_DIR := $(shell basename $(PWD))
# Keep an existing GOPATH, make a private one if it is undefined
GOPATH_DEFAULT := $(PWD)/.go
export GOPATH ?= $(GOPATH_DEFAULT)
PKG := github.com/ZJU-SEL/capstan
DEST := $(GOPATH)/src/$(PKG)

GOFLAGS :=
TAGS :=
LDFLAGS :=

OUTPUT := _output

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

.PHONY: depend
depend: work

.PHONY: build
build: depend
	cd $(DEST)
	$(GO) build $(GOFLAGS) -a -o $(OUTPUT)/capstan ./cmd/capstan

.PHONY: clean
clean:
	rm -rf $(OUTPUT)

.PHONY: install
install: depend
	rm -rf /etc/capstan
	mkdir -p /etc/capstan/prometheus /etc/capstan/grafana/provisioning
	mkdir /etc/capstan/grafana/provisioning/datasources /etc/capstan/grafana/provisioning/dashboards
	cp $(GOPATH)/src/github.com/ZJU-SEL/capstan/grafana-dashboards/* /etc/capstan/grafana/provisioning/dashboards/
	cp $(GOPATH)/src/github.com/ZJU-SEL/capstan/deploy/grafana-datasources.yaml /etc/capstan/grafana/provisioning/datasources/prometheus.yaml
	cp $(GOPATH)/src/github.com/ZJU-SEL/capstan/examples/capstan.conf /etc/capstan/config
	cd $(DEST)
	install -D -m 755 $(OUTPUT)/capstan /usr/local/bin/capstan

.PHONY: uninstall
uninstall:
	rm -f /usr/local/bin/capstan

.PHONY: test
test: test-unit

.PHONY: test-unit
test-unit: depend
test-unit: TAGS += unit
test-unit: test-flags

.PHONY: test-flags
test-flags:
	cd $(DEST) && go test $(GOFLAGS) -tags '$(TAGS)' $(go list ./... | grep -v vendor)

.PHONY: gofmt
gofmt:
	hack/verify-gofmt.sh

.PHONY: lint
lint:
	hack/verify-lint.sh

.PHONY: boiler
boiler:
	hack/verify-boilerplate.sh

.PHONY: install.tools
install.tools:
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install

cover:
	@echo "$@ not yet implemented"

docs:
	@echo "$@ not yet implemented"

godoc:
	@echo "$@ not yet implemented"

# Set up the development environment
.PHONY: env
env:
	@echo "PWD: $(PWD)"
	@echo "BASE_DIR: $(BASE_DIR)"
	@echo "GOPATH: $(GOPATH)"
	@echo "DEST: $(DEST)"
	@echo "PKG: $(PKG)"

work: $(GOPATH) $(DEST)

$(GOPATH):
	mkdir -p $(GOPATH)

$(DEST): $(GOPATH)
	mkdir -p $(shell dirname $(DEST))
	ln -s $(PWD) $(DEST)

shell: work
	cd $(DEST) && $(SHELL) -i