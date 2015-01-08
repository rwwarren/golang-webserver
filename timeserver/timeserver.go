package main

import (
    "fmt"
    "net/http"
    "time"
    "flag"
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
    var port = flag.Int("port", 8080, "help message for flagname")
    flag.Parse()
    http.HandleFunc("/time", handler)
    http.HandleFunc("/", errorer)
    var portString = fmt.Sprintf(":%d", *port)
    fmt.Printf("port string is: %s", portString)
    err:= http.ListenAndServe(portString, nil)
    fmt.Println("Server Failed: %s\n", err)
}

