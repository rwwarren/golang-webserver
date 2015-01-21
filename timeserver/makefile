# Make file for timesrver assignment
# This can get "help", "print", "fmt",
# "install", "run"
# GOLOC is the location of the go executable
# This also takes PORT for the port
# Lastly will take V for version

GOROOT="/usr/apps/go/hg/"
GOLOC="/usr/apps/go/hg/bin/go"
LOGFILE=""

ifdef BREW
	GOLOC=$(shell which go)
endif
ifdef LOG
	LOGFILE="--LogOutput $(LOG)"
	export LOGFILE
endif


all: setup print install run

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
ifdef V
	./timeserver -V
else
  ifndef PORT
	@echo ":8080"
	- ./timeserver $(shell echo $(LOGFILE))
  else
	@echo ":"$(PORT)
	- ./timeserver --port $(PORT) $(shell echo $(LOGFILE))
  endif
endif

