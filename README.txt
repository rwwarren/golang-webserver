this is the readme

tarball command reminder
tar czf timeserver.tar.gz ./css490


For assignment 7:

compile "timeserver" and "authserver"
using go build (so the executable is in the
"command/timserver" and "command/authserver"
directories respectively
Then compile and run "supervisor"
Each component has a readme and a makefile
supervisor (as seen in comments) can load either
a file or os.Args[1] for the json configuration


WARNING: supervisor will not work if the json configuration
is not properly configured.
