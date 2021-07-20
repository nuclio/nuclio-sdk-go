# Copyright 2017 The Nuclio Authors.
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

GOPATH ?= $(shell go env GOPATH)
OS_NAME = $(shell uname)

.PHONY: fmt
fmt:
	gofmt -s -w .

.PHONY: test
test: lint
	go test -v ./...

.PHONY: lint
lint: modules
	@echo Installing linters...
	@test -e $(GOPATH)/bin/impi || \
		(mkdir -p $(GOPATH)/bin && \
		curl -s https://api.github.com/repos/pavius/impi/releases/latest \
			| grep -i "browser_download_url.*impi.*$(OS_NAME)" \
			| cut -d : -f 2,3 \
			| tr -d \" \
			| wget -O $(GOPATH)/bin/impi -qi - \
			&& chmod +x $(GOPATH)/bin/impi)

	@test -e $(GOPATH)/bin/golangci-lint || \
		(curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin v1.41.1)

	@echo Verifying imports...
	chmod +x $(GOPATH)/bin/impi && $(GOPATH)/bin/impi \
		--local github.com/nuclio/nuclio-sdk-go/ \
		--scheme stdLocalThirdParty \
		./...

	@echo Linting...
	$(GOPATH)/bin/golangci-lint run -v
	@echo Done.

.PHONY: modules
modules: ensure-gopath
	@echo Getting go modules
	@go mod download

.PHONY: ensure-gopath
ensure-gopath:
ifndef GOPATH
	$(error GOPATH must be set)
endif
