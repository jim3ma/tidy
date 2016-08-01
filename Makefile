.PHONY: build clean build-docker release release-docker

PROJECT = "tidy"
PATH := $(PATH)
SHELL := /bin/bash
BUILD_DIR = "build"
MAINEXEC = "tidy"

BUILDTAGS=debug

default: all

all: build

deps:
	mkdir -p $(BUILD_DIR)/bin; #\
	#go get -tags '$(BUILDTAGS)' -d -v ./...;

build: deps; \
    CGO_ENABLED=0 go install -tags '$(BUILDTAGS)' . ; \
    cp $${GOPATH}/bin/tidy $(BUILD_DIR)/bin; \
	cp -vfr keys $(BUILD_DIR)/bin/; \
	cp -vfr tidy.yaml $(BUILD_DIR)/bin/;

#build-docker: build; \

build-docker:
	TARGET=$(PROJECT):`date +'%Y-%m-%d'`; \
	docker build -t $${TARGET} .

update-key: build; \
	(cd $(BUILD_DIR)/bin/keys/ && ./key-gen.sh)

release: BUILDTAGS=release
release: build

release-docker: BUILDTAGS=release
#release-docker: build
release-docker: build-docker

update:
	git pull

clean:
	rm -rf build
	rm -f $${GOPATH}/bin/tidy
	go clean -i -r ./...
	git checkout -- .
