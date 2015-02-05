// (C) Ryan Warren 2015
// Authserver
//
// 

package main

import (
    "net/http"
    "os"
    "net"
    "fmt"
)

func init() {
}

func getPath(w http.ResponseWriter, r *http.Request) {
}

func setPath(w http.ResponseWriter, r *http.Request) {
}

func errorer(w http.ResponseWriter, r *http.Request) {
}

func main() {
        ifaces, ipError := net.Interfaces()
        if ipError != nil {
          fmt.Println(ipError)
        }
        for _, i := range ifaces {
            addrs, err := i.Addrs()
            if err != nil {
              fmt.Println(err)
            }
            for _, addr := range addrs {
                //fmt.Println(_)
                fmt.Println(addr)
                switch v := addr.(type) {
                    case *net.IPAddr:
                    // process IP address
                    fmt.Println(v)
                    //fmt.Println(*net.IPAddr)
                }
            }
        }
	http.HandleFunc("/get", getPath)
	http.HandleFunc("/set", setPath)
	http.HandleFunc("/", errorer)
	err := http.ListenAndServe(":9090", nil)
	//err := http.ListenAndServe(portString, nil)
	if err != nil {
		//log.Errorf("Server Failed: %s", err)
		os.Exit(1)
	}
}


