
# Fix this and add port
# Make a little more robust
# GO ?= go
# GOPATH := $(CURDIR)/_vendor:$(GOPATH)
#
# all: build
#
# build:
# 	$(GO) build
# 	godoc -v --http=$(PORT) --index=true
	#$(TESTT)
#all: fmt install run

GOROOT="/usr/apps/go/hg/"
GOHOME=$(HOME)
GOPATH=$(CURDIR)
GOLOC="/usr/apps/go/hg/bin/go"

ifdef BREW
	GOLOC="go"
endif

all: print install run

help:
	@cat readme.txt

print:
	@echo "GOROOT: " $(GOROOT)
	@echo "GOLOC: " $(GOLOC)

fmt:
	@echo "go fmt"
	@$(GOLOC) fmt

install:
	@echo "go build"
	@$(GOLOC) build

run:
	@echo -n "Server location: "
	- @ifconfig | grep -A 2 "eth0" | grep "inet addr" | cut -d: -f2 | awk '{ printf $$1}'

ifndef PORT
	@echo ":8080"
	- ./timeserver
else
	@echo ":"$(PORT)
	- ./timeserver --port $(PORT)
endif

