// (C) Ryan Warren 2015
// Authserver

package main

import (
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
	http.HandleFunc("/get", getPath)
	http.HandleFunc("/set", setPath)
	http.HandleFunc("/", errorer)
	err := http.ListenAndServe(":9090", nil)
	//err := http.ListenAndServe(portString, nil)
	if err != nil {
		log.Errorf("Server Failed: %s", err)
		os.Exit(1)
	}
}


