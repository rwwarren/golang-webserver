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
	"encoding/json"
	"flag"
	"fmt"
	//"io/ioutil"
	//"net/http"
	"os"
	//"strconv"
	//"strings"
	//"sync"
	//"time"
        "io/ioutil"
)

//type mytype []app
type mytype []interface{}
//type mytype []map[string]interface{}

type myapp struct{
//  command interface{}
  //command struct {
  //    command string
  //    args string
  //}
  //command
  output string
  errorString string
}

type commands struct {
      Command string
      Args []string
  }

type app struct {
  command []string
  output string
  errorString string
}

type suchapp struct {
  //Command commands
  //Command commands
  Command []string
  //Command string
  Output string
  Error string
  //command string
  //output string
  //errorString string
}
//type mytype []map[string]string

// Initalizes the 
func init() {
}

func getLoadFile(loadingFile string) []byte {
//func getLoadFile(loadingFile string) string {
  fileBytes, err := ioutil.ReadFile(loadingFile)
  if err != nil {
    log.Criticalf("Failed: %s", err)
    os.Exit(1)
  }
  //TODO leave as byte array?
  return fileBytes
  //return string(fileBytes)
}

func getSupervisionList(loadedFile []byte) []string{
//func getSupervisionList(loadedFile string) []string{

  //fmt.Println(loadedFile)
  //var data []mytype
  var data []suchapp
  //var data []myapp
  //var data mytype
  //var data []app
  //var data mytype
  //keys := make([]mytype,0)

  //err := json.Unmarshal(loadedFile, &keys)
  err := json.Unmarshal(loadedFile, &data)
  //err := json.Unmarshal([]byte(loadedFile), &data)
  if err != nil {
    log.Critical(err)
  }
  //var myslice []mytype
  //DecodeSlicePath(data, &myslice)
  //fmt.Println(data)
  for key, val := range data {
    fmt.Println(key)
    fmt.Println(val)
    //first := val.(map[string]interface{})
    //first := val["commands"].(map[string]interface{})
    //fmt.Println(first)
   // fmt.Println(val["command"])
  }
  return nil
}

// Main function of the supervisor
func main() {
	defer log.Flush()
	logFile := flag.String("log", "logConfig", "This is the logger configuration file")
	portRange := flag.String("port-range", "8080-9090", "This is the port range")
	dumpLoc := flag.String("dumpfile", "backup.bak", "This is the dumpfile")
	loadFile := flag.String("loadfile", "config.json", "This is the dumpfile")
	//loadFile := flag.String("loadfile", "config2.json", "This is the dumpfile")
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
        loadedString := getLoadFile(loadingFile)
        supervisionList := getSupervisionList(loadedString)
        fmt.Println(supervisionList)
        //strings.Replace on {{port}}
}
