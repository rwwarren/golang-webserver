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
        //"os"
        //"os/exec"
        "strings"
	//"bytes"
	//"encoding/json"
	"flag"
	"fmt"
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

func getSupervisionList() []string{

  return nil
}

// Main function of the supervisor
func main() {
	defer log.Flush()
	logFile := flag.String("log", "logConfig", "This is the logger configuration file")
	portRange := flag.String("port-range", "8080-9090", "This is the port range")
	dumpLoc := flag.String("dumpfile", "backup.bak", "This is the dumpfile")
	loadFile := flag.String("loadfile", "config.json", "This is the dumpfile")
        checkout := flag.Int("checkpoint-interval", 2, "This is the checkpoint interval")
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
        ports := *portRange
	portsList := strings.Split(ports, "-")
        dumpfile := *dumpLoc
        checkoutInterval := *checkout
        loadingFile := *loadFile
	log.Infof("port range Flag: %s", ports)
	log.Infof("port range list: %s", portsList)
	log.Infof("dumpfile FLag: %s", dumpfile)
	log.Infof("checkpoint interval Flag: %v", checkoutInterval)
	log.Infof("loading file Flag: %s", loadingFile)
        //supervisionList := getSupervisionList
        //strings.Replace on {{port}}
}
