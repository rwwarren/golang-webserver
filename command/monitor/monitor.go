// (C) Ryan Warren 2015
// Monitor
//
// Monitor it
//
// There are a couple flags for this program:
//--
//--
//--

package main

import (
	log "../../seelog-master/"
	//"../counter"
	"flag"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func collectStats(url string, sampleSec int) {
	for {
		time.Sleep(time.Duration(sampleSec) * time.Second)
		//client := http.Client
		//client := http.Client{
		//	Timeout: (time.Duration(timeout) * time.Millisecond),
		//}
		//response, err := client.Get(url)
		response, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(response)
		}

	}

}

// Main function of the loadgen
func main() {
	defer log.Flush()
	logFile := flag.String("log", "logConfig", "This is the logger configuration file")
	targetList := flag.String("targets", "http://google.com/,http://facebook.com/", "This is the target list of urls")
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
		//  go run()
		go collectStats(url, sampleSec)
	}
	time.Sleep(time.Duration(monitorRuntime) * time.Second)
	//printResults()
	//TODO make the targets: split it on the commas
}
