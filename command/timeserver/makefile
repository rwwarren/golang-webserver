# Make file for timesrver assignment
# This can get "help", "commits", "print",
# "test", "fmt", "install", "run"
# GOLOC is the location of the go executable
# This also takes PORT for the port
# Another thing is LOG and config filename
# for the logfile name (which appends .xml)
# also prepends "etc/" for the file location
# The TEMPLATE flag takes the directory that
# the templates are in
# Lastly will take V for version
    -authhost="http://localhost": This is the authserver default host
      -authport=9090: This is the authserver default port
        -authtimeout-ms=1000: This is the authserver timeout
	  -avg-response-ms=1000: This is the timeserver avg response time
	    -deviation-ms=10: This is the timeserver deviation
	        -max-inflight=0: Max number of inflight requests

GOROOT="/usr/apps/go/hg/"
GOLOC="/usr/apps/go/hg/bin/go"
LOGFILE=""
TEMPLATEFOLDER=""
PORTNUM=""
DEFAULTPORT="8080"

ifdef BREW
	GOLOC=$(shell which go)
endif
ifdef LOG
	LOGFILE="--log $(LOG)"
endif
ifdef TEMPLATE
	TEMPLATEFOLDER="--templates $(TEMPLATE)"
endif
ifdef PORT
	PORTNUM="--port $(PORT)"
	DEFAULTPORT="$(PORT)"
endif


all: print install run

help:
	@cat readme.txt

commits:
	@echo "Total commits on branch: "
	@echo $(shell git rev-list --count HEAD)
print:
	@echo "GOROOT: " $(GOROOT)
	@echo "GOLOC: " $(GOLOC)

fmt:
	@echo "go fmt"
	@$(GOLOC) fmt

test:
	@echo "go test"
	@$(GOLOC) test

install:
	@echo "go build"
	@$(GOLOC) build

run:
	@echo -n "Server location: "
	- @ifconfig | grep -A 2 "eth0" | grep "inet addr" | cut -d: -f2 | awk '{ printf $$1}'
ifdef V
	./timeserver -V
else
	@echo ":"$(DEFAULTPORT)
	- ./timeserver $(shell echo $(PORTNUM)) $(shell echo $(LOGFILE)) $(shell echo $(TEMPLATEFOLDER))
endif

