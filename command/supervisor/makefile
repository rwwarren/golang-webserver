# Make file for timesrver assignment
# This can get "help", "commits", "print",
# "test", "fmt", "install", "run"
# GOLOC is the location of the go executable
# GOROOT is the goo root location
# LOG is logger configuration xml file
# PORTS is the port range for the supervisor
# DUMP is the backup dumpfile for
# LOAD is the loadfile for the json config
# CHECK is the checkpoint interval for the supervisor
# processes command on a mac will show all processes
# add "PORT=:PORTNUM" to search on that port

GOROOT="/usr/apps/go/hg/"
GOLOC="/usr/apps/go/hg/bin/go"
LOGFILE=""
RANGE=""
DUMPFILE=""
LOADFILE=""
CHECKPOINT=""
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
ifdef PORTS
	RANGE="--port-range $(PORTS)"
endif
ifdef DUMP
	DUMPFILE="--dumpfile $(DUMP)"
endif
ifdef LOAD
	LOADFILE="--loadfile $(LOADFILE)"
endif
ifdef CHECK
	CHECKPOINT="--checkpoint-interval $(CHECKPOINT)"
endif

all: print install run

help:
	@cat readme.txt

commits:
	@echo "Total commits on branch: "
	@echo $(shell git rev-list --count HEAD)

processes:
	lsof -i $(PORT)

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
	- ./supervisor $(shell echo $(LOGFILE)) $(shell echo $(RANGE)) $(shell echo $(DUMPFILE)) $(shell echo $(LOADFILE)) $(shell echo $(CHECKPOINT))

