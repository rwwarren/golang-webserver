// (C) Ryan Warren 2015
// Supervisor
//
// Supervisor does....
//
// There are a couple flags for this program:
// --log is the logger configuration file
// --port-range
// --checkpoint-interval interval between checking services
// --dumpfile sdf

package main

import (
	log "../../seelog-master/"
        "os"
        "os/exec"
        "strings"
	//"bytes"
	//"encoding/json"
	//"flag"
	//"fmt"
	//"io/ioutil"
	//"net/http"
	//"os"
	//"strconv"
	//"strings"
	//"sync"
	//"time"
)


// Initalizes the 
func init() {
}

// Main function of the supervisor
func main() {
	defer log.Flush()
	logFile := flag.String("log", "logConfig", "This is the logger configuration file")
	//targetList := flag.String("targets", "http://localhost:8080,http://localhost:9090", "This is the target list of urls")
	//sampleInterval := flag.Int("sample-interval-sec", 2, "This is the sample request rate")
	//runtime := flag.Int("runtime-sec", 20, "This is the monitor runtime")
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
	//targets := *targetList
	//sampleSec := *sampleInterval
	//monitorRuntime := *runtime
	//targetsList := strings.Split(targets, ",")
	//log.Infof("target list Flag: %s", targets)
	//log.Infof("sample-interval-sec Flag: %v", sampleSec)
	//log.Infof("runtime Flag: %v", monitorRuntime)
}
