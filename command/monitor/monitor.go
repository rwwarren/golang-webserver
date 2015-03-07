// (C) Ryan Warren 2015
// Monitor
//
// Monitor monitors the stats of a list of targets
// it checks each target at a specific interval to
// collect useage statistics for the user to view
//
// There are a couple flags for this program:
// --log is the logger configuration file
// --targets is a comma separated list of urls to monitor
// --sample-interval-sec interval between samples
// --runtime-sec runtime of the monitor before it is done

package main

import (
	log "../../seelog-master/"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Stores the sample we have taken, along with the
// time value and average
type Sample struct {
	Time    string
	Value   map[string]interface{}
	Average map[string]interface{}
}

// Stores the target response information
var concurrentMap struct {
	sync.RWMutex
	target map[string][]Sample
}

// Initalizes the authserver with the important user storage things
func init() {
	concurrentMap = struct {
		sync.RWMutex
		target map[string][]Sample
	}{target: make(map[string][]Sample)}
}

// Prints the results of all the the targets
func printResults(interval int) {
	concurrentMap.RLock()
	sample := concurrentMap.target
	concurrentMap.RUnlock()
	fmt.Println()
	fmt.Println()
	fmt.Println("Results (averages):")
	fmt.Println()
	fmt.Println()
	for keys, values := range sample {
		fmt.Println(keys)
		for i := 1; i < len(values); i++ {
			jsonString, err := json.Marshal(values[i].Average)
			if err != nil {
				fmt.Println(err)
			}
			var out bytes.Buffer
			json.Indent(&out, jsonString, "", "\t")
			out.WriteTo(os.Stdout)
			fmt.Println()

		}
		fmt.Println()
	}
}

// Gets ad saves the json from the current target
func request(url string, sampleSec int, pastJson map[string]interface{}) map[string]interface{} {
	var jsonResults map[string]interface{}
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	} else {
		body, _ := ioutil.ReadAll(response.Body)
		respBody := string(body)
		log.Infof("Requested url: %s. Response body: %s", url, respBody)
		concurrentMap.Lock()
		sampleSlice := concurrentMap.target[url]
		err := json.Unmarshal([]byte(respBody), &jsonResults)
		if err != nil {
			log.Errorf("%v", err)
		}
		if pastJson != nil {
			for key, value := range jsonResults {
				firstValueStr := value.(string)
				secondValueStr := pastJson[key].(string)
				firstValue, firstErr := strconv.Atoi(firstValueStr)
				secondValue, secondErr := strconv.Atoi(secondValueStr)
				if firstErr != nil || secondErr != nil {
					fmt.Println(firstErr)
					fmt.Println(secondErr)
				}
				avg := ((firstValue - secondValue) / sampleSec)
				pastJson[key] = avg
			}
		}
		newSample := Sample{fmt.Sprintf("%v", time.Now()), jsonResults, pastJson}
		sampleSlice = append(sampleSlice, newSample)
		concurrentMap.target[url] = sampleSlice
		concurrentMap.Unlock()
	}
	return jsonResults
}

// Collects the stats for the current target evert sampleSec seconds
func collectStats(url string, sampleSec int) {
	firstResult := request(url, sampleSec, nil)
	for {
		time.Sleep(time.Duration(sampleSec) * time.Second)
		firstResult = request(url, sampleSec, firstResult)
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
	printResults(sampleSec)
}
