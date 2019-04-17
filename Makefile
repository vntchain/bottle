.PHONY: bottle 

GOBIN = $(shell pwd)/build/bin
GO ?= latest
ARCH=$(shell go env GOARCH)
MARCH=$(shell go env GOOS)-$(shell go env GOARCH)

all: bottle

bottle: 
	scripts/env.sh  scripts/build.sh 
	
bottle-docker:
	docker build -t vntchain/bottle:0.6.0 ./docker/ubuntu

clean:
	rm -fr build/