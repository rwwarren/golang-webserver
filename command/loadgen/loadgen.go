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
	"sort"
	"sync"
	"time"
        "os"
)

// Stores the cookie information
var concurrentMap struct {
	sync.RWMutex
	cookieMap map[string]int
}
var c = counter.New()
var convert   = map[int]string{
                1: "100s",
                2: "200s",
                3: "300s",
                4: "400s",
                5: "500s",
        }

// Initalizes the loadgen with the important user storage things
func init() {
	concurrentMap = struct {
		sync.RWMutex
		cookieMap map[string]int
	}{cookieMap: make(map[string]int)}
	concurrentMap.Lock()
	concurrentMap.cookieMap["Total"] = 0
	concurrentMap.cookieMap["100s"] = 0
	concurrentMap.cookieMap["200s"] = 0
	concurrentMap.cookieMap["300s"] = 0
	concurrentMap.cookieMap["400s"] = 0
	concurrentMap.cookieMap["500s"] = 0
	concurrentMap.cookieMap["Errors"] = 0
	concurrentMap.Unlock()
}

// Gets the status code range from the statusCode
func getStatusCentury(statusCode int) string {
	keyCode := (statusCode / 100) * 100
        if keyCode > 500 || keyCode < 100 {
          return fmt.Sprint("Errors")
        }
	return fmt.Sprintf("%vs", keyCode)
}

func getUrl(testUrl string, timeout int) {
	getUrl := fmt.Sprintf("%s", testUrl)
	timeoutms := time.Duration(time.Duration(timeout) * time.Millisecond)
	client := http.Client{
		Timeout: timeoutms,
	}
	resp, err := client.Get(getUrl)
	if err != nil {
		log.Criticalf("Error getting authserver: %s", err)
		//this is a timout / error
		//return ""
		concurrentMap.Lock()
		currentCount := concurrentMap.cookieMap["Errors"]
		currentCount++
		concurrentMap.cookieMap["Errors"] = currentCount
		totalCount := concurrentMap.cookieMap["total"]
		totalCount++
		concurrentMap.cookieMap["total"] = totalCount
		concurrentMap.Unlock()
	} else {
		log.Infof("Response from the authserver: %s", resp)
		defer resp.Body.Close()
		status := resp.StatusCode
		mapStatus := fmt.Sprintf("%v", getStatusCentury(status))
		concurrentMap.Lock()
		currentCount := concurrentMap.cookieMap[mapStatus]
		currentCount++
		concurrentMap.cookieMap[mapStatus] = currentCount
		totalCount := concurrentMap.cookieMap["Total"]
		totalCount++
		concurrentMap.cookieMap["Total"] = totalCount
		concurrentMap.Unlock()
	}
}

func printMap(runtime int) {
	//time.Sleep(time.Duration(runtime) * time.Second)
	concurrentMap.RLock()
	//TODO print the map

	keys := []string{}
	for k, _ := range concurrentMap.cookieMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	//sort.Ints(keys)
	for _, k := range keys {
		//log.Infof("%v: %v", k, concurrentMap.cookieMap[k])
		//fmt.Printf("%v: %v \n", k, concurrentMap.cookieMap[k])
		fmt.Printf("%v: %v \n", k, c.Get(k))
		//fmt.Printf("%v: %v \n", k, convert[k])
	}

	//for key, val := range concurrentMap.cookieMap {
	//    log.Debugf("%v: %v", key, val)
	//    fmt.Printf("%v: %v\n", key, val)
	//  }
	concurrentMap.RUnlock()
}
func request(timeout int, testUrl string) {
        c.Incr("Total", 1)
        client := http.Client{
                Timeout: (time.Duration(timeout) * time.Millisecond),
        }
        response, err := client.Get(testUrl)
        if err != nil {
                fmt.Printf("err: %s \n", err)
                c.Incr("Errors", 1)
                return
        }

        //if response.StatusCode < 200 {
        //        c.Incr("100s", 1)
        //} else if response.StatusCode < 300 {
        //        c.Incr("200s", 1)
        //}

        //c.Incr(fmt.Sprintf("%ds", (response.StatusCode/100)*100), 1)

        //key, ok := convert[response.StatusCode/100]
		//concurrentMap.Lock()
		//currentCount := concurrentMap.cookieMap["Errors"]
		//concurrentMap.cookieMap["Errors"] = currentCount
		//totalCount := concurrentMap.cookieMap["total"]
		//concurrentMap.cookieMap["total"] = totalCount
                key := fmt.Sprintf("%v", getStatusCentury(response.StatusCode))
                //key := fmt.Sprintf("%v", getStatusCentury(status))
        //key, ok := concurrentMap.cookieMap[fmt.Sprintf("%v", getStatusCentury(status))]
		//concurrentMap.Unlock()
        //key, ok := convert[mapStatus := fmt.Sprintf("%v", getStatusCentury(status))]
        //if !ok {
        //        key = "errors"
        //}
        c.Incr(key, 1)
}

func load(testUrl string, reqRate int, burstRate int, timeout int, runtime int) {
        //timeoutTick := time.Tick(time.Duration(timeout) * time.Second)
        timeoutTick := time.Tick(time.Duration(runtime) * time.Second)
// maybe runtime??????????
        //timeoutTick := time.Tick(runtime * time.Second)
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
	//
	load(testUrl, reqRate, burstRate, timeout, runtime)
	//getUrl(testUrl, timeout)
	//go printMap(runtime)
        //time.Sleep(time.Duration(runtime) * time.Second)
	time.Sleep(time.Duration(2 * runtime) * time.Second)
	printMap(runtime)
        os.Exit(0)
}
