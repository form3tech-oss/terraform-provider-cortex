include Makefile.deps

TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=registry.terraform.io
NAMESPACE=form3tech-oss
NAME=cortex
BINARY=terraform-provider-${NAME}
VERSION=0.0.4
OS_ARCH=linux_amd64
PATH := $(PATH):$(PWD)/bin
SHELL := /bin/bash

default: install

.PHONY: build
build:
	go build -o ${BINARY}

.PHONY: release
release: deps
	goreleaser

.PHONY: install
install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

.PHONY: test
test: 
	go test -v ./...

.PHONY: testacc
testacc: deps
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m   

.PHONY: clean
clean:
	rm -f examples/terraform.tfstate
	rm -f examples/terraform.tfstate.backup

.PHONY: lint
lint: vet tflint tffmtcheck

.PHONY: vet
vet:
	go vet ./...

.PHONY: tflint
tflint: deps
	find ./examples/ -type d -exec tflint \{\} \;

.PHONY: tffmtcheck
tffmtcheck: deps
	terraform fmt -check -recursive ./examples/

.PHONY: fmt
fmt: deps
	go fmt ./...
	terraform fmt -recursive ./examples/

.PHONY: docs
docs: deps
	tfplugindocs generate

dev.tfrc:
	echo 'provider_installation {' >> dev.tfrc
	echo '  dev_overrides {' >> dev.tfrc
	echo '    "form3tech-oss/cortex" = "$(CURDIR)"' >> dev.tfrc
	echo '  }' >> dev.tfrc
	echo '  direct {}' >> dev.tfrc
	echo '}' >> dev.tfrc

.PHONY: cortex-up
cortex-up:
	docker-compose up -d

.PHONY: cortex-down
cortex-down:
	docker-compose down
