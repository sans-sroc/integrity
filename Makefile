BRANCH := $(shell git rev-parse --symbolic-full-name --abbrev-ref HEAD)
SUMMARY := $(shell bash .ci/version)
VERSION := $(shell cat VERSION)
NAME := $(shell basename `pwd`)
MODULE := $(shell cat go.mod | head -n1 | cut -f2 -d' ')

.PHONY: build release vendor

vendor:
	go mod vendor

build: vendor
	go build -ldflags "-X $(MODULE)/pkg/common.SUMMARY=$(SUMMARY) -X $(MODULE)/pkg/common.BRANCH=$(BRANCH) -X $(MODULE)/pkg/common.VERSION=$(VERSION)" -o $(NAME)

release: vendor
	mkdir -p release
	go build -mod=vendor -ldflags "-X $(MODULE)/pkg/common.SUMMARY=$(SUMMARY) -X $(MODULE)/pkg/common.BRANCH=$(BRANCH) -X $(MODULE)/pkg/common.VERSION=$(VERSION)" -o release/$(NAME) .

run-%: vendor
	go run -mod=vendor -ldflags "-X $(MODULE)/pkg/common.SUMMARY=$(SUMMARY) -X $(MODULE)/pkg/common.BRANCH=$(BRANCH) -X $(MODULE)/pkg/common.VERSION=$(VERSION)" main.go $*
