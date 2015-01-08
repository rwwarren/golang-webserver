
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

HOME="/usr/apps/go/hg/"
GOHOME=$(HOME)
GOPATH=$(CURDIR)
GOROOT="/usr/apps/go/hg/bin/go"

ifdef BREW
	GOROOT="go"
endif

all: print install run

help:
	@cat readme.txt

print:
	@echo "GOROOT: " $(GOROOT)

fmt:
	@echo "go fmt"
	@$(GOROOT) fmt

install:
	@echo "go build"
	@$(GOROOT) build

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

