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
        "../counter"
	"flag"
	"fmt"
	"net/http"
	"time"
        "os"
)

// Counter
var c = counter.New()

// Keys for tracking
var keys = []string {
  "100s",
  "200s",
  "300s",
  "400s",
  "500s",
  "Errors",
  "Total",
  }

// Gets the status code range from the statusCode
func getStatusCentury(statusCode int) string {
	keyCode := (statusCode / 100) * 100
        if keyCode > 500 || keyCode < 100 {
          return fmt.Sprint("Errors")
        }
	return fmt.Sprintf("%vs", keyCode)
}

// Prints the output information
// about all the requests
func printMap(runtime int) {
	for _, k := range keys {
		fmt.Printf("%s: %v \n", k, c.Get(k))
	}
}

// Sends request to server and tracks status code
func request(timeout int, testUrl string) {
        c.Incr("Total", 1)
        client := http.Client{
                Timeout: (time.Duration(timeout) * time.Millisecond),
        }
        response, err := client.Get(testUrl)
        if err != nil {
                c.Incr("Errors", 1)
                return
        }
                key := fmt.Sprintf("%v", getStatusCentury(response.StatusCode))
        c.Incr(key, 1)
}

// Creates all the bursts and fires off requests
func load(testUrl string, reqRate int, burstRate int, timeout int, runtime int) {
        timeoutTick := time.Tick(time.Duration(runtime) * time.Second)
        interval := time.Duration((1000000*burstRate)/reqRate) * time.Microsecond
        period := time.Tick(interval)
        for {
                // fire off burst
                for i := 0; i < burstRate; i++ {
                        go request(timeout, testUrl)
                }
                // wait for next tick
                <-period
                select {
                case <-timeoutTick:
                        return
                default:
                }
        }
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
	load(testUrl, reqRate, burstRate, timeout, runtime)
        //Added a second to let a little extra collection happen
	time.Sleep(time.Duration(runtime + 1) * time.Second)
        fmt.Println()
	printMap(runtime)
        os.Exit(0)
}
