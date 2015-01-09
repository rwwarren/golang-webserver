//Testing
// find out how to see this
//
//
package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
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

func main() {
	port := flag.Int("port", 8080, "Set the server port, default port: 8080")
	version := flag.Bool("V", false, "shows the version?")
	//var port = flag.Int("port", 8080, "Set the server port, default port: 8080")
	flag.Parse()
        if(*version){
	  fmt.Println("Assignment Version: 1")
          return
        }
	http.HandleFunc("/time", handler)
	http.HandleFunc("/", errorer)
	var portString = fmt.Sprintf(":%d", *port)
	err := http.ListenAndServe(portString, nil)
	fmt.Printf("Server Failed: %s\n", err)
}
