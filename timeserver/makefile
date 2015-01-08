
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

HOME="/usr/apps/go/hg/"
GOHOME=$(HOME)
GOPATH=$(CURDIR)

all: fmt install

fmt:
	go fmt
install:
#run:
	echo $(PORT)
	go build
ifndef PORT
	./timeserver
else
	./timeserver --port $(PORT)
endif
