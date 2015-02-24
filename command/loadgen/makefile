# Make file for timesrver assignment
# This can get "help", "commits", "print",
# "test", "fmt", "install", "run"
# GOLOC is the location of the go executable
# GOROOT is the goo root location
# LOG is logger configuration xml file
# TESTSERVER is the url for the loadgen
# to test. RATE is the rate of requests of
# loadgen. BURST is the number of concurrent
# requests. TIMEOUTTIME is the time to wait
# for a request to come back. Finally
# TIMEOUTTIME is the max time to
# process each request.

GOROOT="/usr/apps/go/hg/"
GOLOC="/usr/apps/go/hg/bin/go"
LOGFILE=""
TESTURL=""
REQRATE=""
BURSTRATE=""
TIMEOUT=""
RUNTIME=""
GOTYPE := $(shell { command -v $(GOLOC) ; } 2>/dev/null)
ifneq ("$(GOTYPE)", $(GOLOC))
	GOLOC="go"
endif

ifdef BREW
	GOLOC=$(shell which go)
endif
ifdef LOG
	LOGFILE="--log $(LOG)"
endif
ifdef TESTSERVER
	TESTURL="--url $(TESTSERVER)"
endif
ifdef RATE
	REQRATE="--rate $(RATE)"
endif
ifdef BURST
	BURSTRATE="--burst $(BURST)"
endif
ifdef TIMEOUTTIME
	TIMEOUT="--timeout-ms $(TIMEOUTTIME)"
endif
ifdef RUN
	RUNTIME="--runtime $(RUN)"
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
	- ./loadgen $(shell echo $(LOGFILE)) $(shell echo $(TESTURL)) $(shell echo $(REQRATE)) $(shell echo $(BURSTRATE)) $(shell echo $(TIMEOUT)) $(shell echo $(RUNTIME))

