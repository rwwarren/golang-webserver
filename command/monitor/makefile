# Make file for timesrver assignment
# This can get "help", "commits", "print",
# "test", "fmt", "install", "run"
# GOLOC is the location of the go executable
# GOROOT is the goo root location
# LOG is logger configuration xml file
# TARGET is the comma seperated list of urls
# for the monitor to test. RATE is the rate of
# sampling. RUN is the monitor runtime

GOROOT="/usr/apps/go/hg/"
GOLOC="/usr/apps/go/hg/bin/go"
LOGFILE=""
TARGETS=""
RATE=""
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
ifdef TARGET
	TARGETS="--targets $(TARGET)"
endif
ifdef RATE
	SAMPLE="--sample-interval-sec $(RATE)"
endif
ifdef RUN
	RUNTIME="--runtime-sec $(RUN)"
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
	- ./monitor $(shell echo $(LOGFILE)) $(shell echo $(TARGET)) $(shell echo $(RATE)) $(shell echo $(RUNTIME))

