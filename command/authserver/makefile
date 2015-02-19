# Make file for timesrver assignment
# This can get "help", "commits", "print",
# "test", "fmt", "install", "run"
# GOLOC is the location of the go executable
# This also takes PORT for the port
# Another thing is LOG and config filename
# for the logfile name (which appends .xml)
# also prepends "etc/" for the file location
# The DUMPFILE flag takes the file that
# to backup to
# Lastly will take CHECKINT flag is the checkpoint
# interval for backups


GOROOT="/usr/apps/go/hg/"
GOLOC="/usr/apps/go/hg/bin/go"
LOGFILE=""
DUMPFILENAME=""
PORTNUM=""
DEFAULTPORT="9090"
CHECKPOINT=""

ifdef BREW
	GOLOC=$(shell which go)
endif
ifdef LOG
	LOGFILE="--log $(LOG)"
endif
ifdef DUMPFILE
	DUMPFILENAME="--dumpfile $(DUMPFILE)"
endif
ifdef PORT
	PORTNUM="--port $(PORT)"
	DEFAULTPORT="$(PORT)"
endif
ifdef CHECKINT
	CHECKPOINT="--checkpoint-interval $(CHECKINT)"
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
	@echo ":"$(DEFAULTPORT)
	- ./authserver $(shell echo $(PORTNUM)) $(shell echo $(LOGFILE)) $(shell echo $(TEMPLATEFOLDER))

