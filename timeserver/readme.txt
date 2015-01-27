
  490 Assignment 3 By Ryan Warren
  Here is the README and instructions for it

  Usage: make [options] [variables]

  Examples:   make
              make all
              make run PORT=3030
              make run V=1
              make run BREW=1
              make run LOG=asdf
              make run TEMPLATE=asdf

  Options:
    help      Displays this, the help message
    all       Print, install, then run
    commits   Prints the amount of commits on the branch
    print     Prints the GOROOT
    fmt       Runs go fmt
    test      Runs go test
    install   Runs "go build" to compile the go program
    run       Runs the compiled go build

  Variables:
    BREW      Set to anything if "go" is set up properaly in the $PATH
    GOLOC     If unset, then GOLOC will default to "/usr/apps/go/hg/bin/go"
              otherwise set it to the location where go is installed
              It is the location where go is installed on the computer
    GOROOT    If unset, then GOROOT will default to "/usr/apps/go/hg/"
              otherwise set it to the location where go is installed
    PORT      Set the port for the server to run on
    LOG       Set the log file config name
    TEMPLATE  Set the template directory location
    V         Will display the server version and quit


