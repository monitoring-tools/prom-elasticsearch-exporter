TARGET              ?= prom-elasticsearch-exporter

GOPATH              := $(lastword $(subst :, ,$(GOPATH)))
GIT_SUMMARY         := $(shell git describe --tags --dirty --always)
GIT_BRANCH          := $(shell git rev-parse --abbrev-ref HEAD)
GO_VERSION          := $(shell go version)
LDFLAGS             := -ldflags "-X 'main.version=$(GIT_SUMMARY)' -X 'main.goVersion=$(GO_VERSION)' -X 'main.gitBranch=$(GIT_BRANCH)'"
DOCKER_BUILD_ARGS   ?=

.PHONY: all build test docker

all: test build

test:
	@echo ">> running tests"
	@go test $(shell go list ./...)

build: $(TARGET)

$(TARGET):
	@echo ">> building binary..."
	@echo ">> GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) -o $(TARGET)"
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) -o $(TARGET)

docker: GOOS="linux" GOARCH="amd64"
docker: DOCKER_IMAGE_NAME ?= "monitoring-tools/prom-elasticsearch-exporter:$(GIT_SUMMARY)"
docker: Dockerfile build
	@echo ">> building docker image"
	@docker build -t $(DOCKER_IMAGE_NAME) $(DOCKER_BUILD_ARGS) .
