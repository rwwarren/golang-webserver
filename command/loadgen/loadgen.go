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

func runsomehitng(){

}

// Main function of the loadgen
func main() {
	defer log.Flush()
	logFile := flag.String("log", "logConfig", "This is the logger configuration file")
	testServerUrl := flag.String("url", "http://localhost:8080/time", "This is the test server url")
	requestRate := flag.Int("rate", 200, "This is the request rate")
	burstRequest := flag.Int("burst", 20, "Number of concurrent requests to issue")
	timeoutTime := flag.Int("timeout-ms", 1000, "Max time to wait for response")
	totalRuntime := flag.Int("runtime", 10, "Number of seconds to process")
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
        testUrl := *testServerUrl
        reqRate := *requestRate
        burstRate := *burstRequest
        timeout := *timeoutTime
        runtime := *totalRuntime
        log.Infof("url Flag: %s", testUrl)
        log.Infof("rate Flag: %v", reqRate)
        log.Infof("burst Flag: %v", burstRate)
        log.Infof("timeout-ms Flag: %v", timeout)
        log.Infof("runtime Flag: %v", runtime)
}
