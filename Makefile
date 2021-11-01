TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=registry.terraform.io
NAMESPACE=form3tech-oss
NAME=cortex
BINARY=terraform-provider-${NAME}
VERSION=0.0.5
OS_ARCH=linux_amd64
PROJECT_NAME ?= $(CURDIR)

default: install

build:
	go build -o ${BINARY}

release:
	goreleaser

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

test: 
	go test -i $(TEST) || exit 1                                                   
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4                    

testacc: 
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m   

dev.tfrc:
	echo 'provider_installation {' >> dev.tfrc
	echo '  dev_overrides {' >> dev.tfrc
	echo '    "form3tech-oss/cortex" = "$(CURDIR)"' >> dev.tfrc
	echo '  }' >> dev.tfrc
	echo '  direct {}' >> dev.tfrc
	echo '}' >> dev.tfrc

.PHONY: deps
deps: bin/goreleaser

bin/goreleaser:
	mkdir -p bin
	wget https://github.com/goreleaser/goreleaser/releases/download/v0.164.0/goreleaser_Linux_x86_64.tar.gz
	tar -C bin -xzf goreleaser*tar.gz goreleaser
	rm goreleaser*tar.gz*

