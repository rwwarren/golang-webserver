
  490 Assignment 7 By Ryan Warren
  Here is the README and instructions for it

  Usage: make [options] [variables]

  Examples:   make
              make all
              make run BREW=1
              make run LOG=asdf
              make run PORTS=8080-9090
              make run LOAD=fdas
              make run DUMP=asdf
              make run CHECK=2

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
    PORTS     Set the port range for the supervisor
    LOG       Set the log file config name. Set the LOG to "assignmentLog"
              for the one based on the project specs
    DUMP      Set the backup file location
    LOAD      Set loadfile for the supervisor
    CHECK     Set the check interval

