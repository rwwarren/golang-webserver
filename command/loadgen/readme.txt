
  490 Assignment 5 By Ryan Warren
  Here is the README and instructions for it

  Usage: make [options] [variables]

  Examples:     make
                make all
                make run TESTSERVER=http://localhost:8080/time
                make run RATE=213
                make run BREW=1
                make run LOG=asdf
                make run BURST=10
                make run TIMEOUTTIME=10
                make run RUN=20

  Options:
    help        Displays this, the help message
    all         Print, install, then run
    commits     Prints the amount of commits on the branch
    print       Prints the GOROOT
    fmt         Runs go fmt
    test        Runs go test
    install     Runs "go build" to compile the go program
    run         Runs the compiled go build

  Variables:
    BREW        Set to anything if "go" is set up properaly in the $PATH
    GOLOC       If unset, then GOLOC will default to "/usr/apps/go/hg/bin/go"
                otherwise set it to the location where go is installed
                It is the location where go is installed on the computer
    GOROOT      If unset, then GOROOT will default to "/usr/apps/go/hg/"
                otherwise set it to the location where go is installed
    LOG         Log configuration file for the logger
    TESTSERVER  Test server url for the loadgen to test
    RATE        Set the rate of requests per second
    BURST       Set the number of concurrent requests
    TIMEOUTTIME Set the max wait time for a request
    RUN         Set the number of seconds to process


