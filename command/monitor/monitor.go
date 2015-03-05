// (C) Ryan Warren 2015
// Monitor
//
// Monitor it
//
// There are a couple flags for this program:
// --log is the logger configuration file
// --targets is a comma separated list of urls to monitor
// --sample-interval-sec interval between samples
// --runtime-sec runtime of the monitor before it is done

package main

import (
	log "../../seelog-master/"
	"flag"
	"fmt"
	"net/http"
	"strings"
	"time"
	"io/ioutil"
	"sync"
)

// Stores the cookie information
var concurrentMap struct {
	sync.RWMutex
	cookieMap map[string]string
}

// Initalizes the authserver with the important user storage things
func init() {
  // map[string] []Sample
	concurrentMap = struct {
		sync.RWMutex
		cookieMap map[string]string
	}{cookieMap: make(map[string]string)}
}

func printResults(){
	//concurrentMap.RLock()
	//name := concurrentMap.cookieMap[formCookie]
	//concurrentMap.RUnlock()
}

func collectStats(url string, sampleSec int) {
	for {
		time.Sleep(time.Duration(sampleSec) * time.Second)
		response, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
		} else {
			//fmt.Println(response)
                        fmt.Println(url)
                        body, _ := ioutil.ReadAll(response.Body)
                        //body, err := ioutil.ReadAll(response.Body)
	                respBody := string(body)
			fmt.Println(respBody)
			//fmt.Println(response.Body)
		}
	}
}

// Main function of the loadgen
func main() {
	defer log.Flush()
	logFile := flag.String("log", "logConfig", "This is the logger configuration file")
	targetList := flag.String("targets", "http://localhost:8080,http://localhost:9090", "This is the target list of urls")
	sampleInterval := flag.Int("sample-interval-sec", 2, "This is the sample request rate")
	runtime := flag.Int("runtime-sec", 20, "This is the monitor runtime")
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
	targets := *targetList
	sampleSec := *sampleInterval
	monitorRuntime := *runtime
	targetsList := strings.Split(targets, ",")
	log.Infof("target list Flag: %s", targets)
	log.Infof("sample-interval-sec Flag: %v", sampleSec)
	log.Infof("runtime Flag: %v", monitorRuntime)
	for _, url := range targetsList {
                monitorUrl := fmt.Sprintf("%s/monitor", url)
		go collectStats(monitorUrl, sampleSec)
	}
	time.Sleep(time.Duration(monitorRuntime) * time.Second)
	printResults()
	//TODO make the targets: split it on the commas
}
