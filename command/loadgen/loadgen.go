// (C) Ryan Warren 2015
// Loadgen
//
// Authserver. This tracks all the user logged in information
// it stores infomation about the user to make sure that another
// server can see if the user is logged in or not
//
// There are a couple flags for this program:
//--rate: average rate of requests (per second)
//--burst: number of concurrent requests to issue
//--timeout-ms: max time to wait for response
//--runtime: number of seconds to process
//--url: URL to sample

package main

import (
	log "../../seelog-master/"
	"flag"
        "sync"
        "fmt"
)


// Stores the cookie information
var concurrentMap struct {
	sync.RWMutex
	cookieMap map[string]string
}

// Initalizes the loadgen with the important user storage things
func init() {
	concurrentMap = struct {
		sync.RWMutex
		cookieMap map[string]string
	}{cookieMap: make(map[string]string)}
}


// Main function of the loadgen
func main() {
	defer log.Flush()
	logFile := flag.String("log", "logConfig", "This is the logger configuration file")
	flag.Parse()
	logFileName := fmt.Sprintf("etc/%s.xml", *logFile)
	logger, logError := log.LoggerFromConfigAsFile(logFileName)
	if logError != nil {
		fmt.Printf("Log instantiation error: %s", logError)
	}
	log.ReplaceLogger(logger)
	log.Debug("Logger intitalized")
	log.Trace("Testing trace")
	log.Debug("Testing debug")
	log.Info("Testing info")
	log.Warn("Testing warn")
	log.Error("Testing error")
	log.Critical("Testing critical")
}
