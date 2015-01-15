// (C) Ryan Warren 2015
// Timeserver
//
// This is a timeserver that shows the time for the current timezone
// Hours, Minutes, and Seconds on the "/time" url.
// All other urls will return with a 404.
//
// Building and running with go can be done with "go build"
// then "./timeserver". This will show up on port 8080 by default.
// If the port is in use there will be an error.
//
// There are command line options.
// "--port PORT_NUMBER" will change the port.
// "-V" will show the version and then quit.

package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"
        //"os/exec"
)

// Handles the timeserver which shows the current time
// for the local timezone
func timeHandler(w http.ResponseWriter, r *http.Request) {
	const layout = "3:04:05 PM"
	fmt.Fprintf(w, `<html><head><style>
          p {font-size: xx-large}
          span.time {color: red}
          </style>
          </head>
          <body>
          <p>The time is now <span class="time">%s</span>.</p>
          </body>
          </html>`, time.Now().Local().Format(layout))
}

// Handles errors for when the page is not found
func errorer(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	fmt.Fprintf(w, `<html><head><style>
          p {font-size: xx-large}
          span.time {color: red}
          </style>
          </head>
          <body>
          <p>These are not the URLs you're looking for.</p>
          </body>
          </html>`)
}

func loginForm(w http.ResponseWriter, r *http.Request) {
        //cookie := http.Cookie{"test", "tcookie", "/", "www.sliceone.com", expire, expire.Format(time.UnixDate), 86400, true, true, "test=tcookie", []string{"test=tcookie"}}
        cookie := &http.Cookie{Name:"name", Value:"ryan", Expires:time.Now().Add(356*24*time.Hour), HttpOnly:true}
        //r.AddCookie(&cookie)
        http.SetCookie(w, cookie)
	fmt.Fprintf(w, `<html>
          <body>
          <form action="login">
            What is your name, Earthling?
            <input type="text" name="name" size="50">
            <input type="submit">
          </form>
          </p>
          </body>
          </html>`)
}

func logoutPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<html>
          <head>
          <META http-equiv="refresh" content="10;URL=/">
          <body>
          <p>Good-bye.</p>
          </body>
          </html>`)
}

// Main handler that runs the server on the port or shows the version of the server
func main() {
	port := flag.Int("port", 8080, "Set the server port, default port: 8080")
	version := flag.Bool("V", false, "Shows the version of the timeserver")
	flag.Parse()
	if *version {
		fmt.Println("Assignment Version: 1")
		return
	}
	http.HandleFunc("/time", timeHandler)
	http.HandleFunc("/", loginForm)
	http.HandleFunc("/logout", logoutPage)
	//http.HandleFunc("/", errorer)
	var portString = fmt.Sprintf(":%d", *port)
	err := http.ListenAndServe(portString, nil)
	fmt.Printf("Server Failed: %s\n", err)
}
