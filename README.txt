This is a golang webserver built for a class.

There are 6 components of the server:
  1. Timeserver: This is the frontend server that has a couple different pages, including one that shows you the current time.
  2. Authserver: This connects to the timeserver to track logging in and out of users.
  3. Counter: This is a counter that allows concorrent access
  4. Loadgen: This creates a load on the timeserver and authserver
  5. Monitor: This is a monitor that checks servers at an interval for a runtime
  6. Supervisor: This makes starts and tracks the status of servers to entire there is no downtime






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

Final Note:
If you run supervisor and hit "COMMAND + c" it will
kill all the processes in the process group (all
controlled servers) and if you hit "Q + ENTER" it
will only kill the supervisor and the other processes
will survive
