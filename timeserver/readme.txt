
  490 Assignment 1 By Ryan Warren
  Here is the README and instructions for it

  Usage: make [options] [variables]

  Examples:   make
              make all
              make run PORT=3030
              make run V=1
              make run BREW=1
              make run LOG=asdf

  Options:
    help      Displays this, the help message
    all       Print, install, then run
    print     Prints the GOROOT
    install   Runs "go build" to compile the go program
    run       Runs the compiled go build
    test      Runs go test

  Variables:
    BREW      Set to anything if go is set up properaly in the $PATH
    GOLOC     If unset, then GOLOC will default to "/usr/apps/go/hg/bin/go"
              otherwise set it to the location where go is installed
              It is the location where go is installed on the computer
    GOROOT    If unset, then GOROOT will default to "/usr/apps/go/hg/"
              otherwise set it to the location where go is installed
    PORT      Set the port for the server to run on
    LOG       Set the log file config name
    V         Will display the server version and quit


