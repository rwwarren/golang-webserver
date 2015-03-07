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
        "encoding/json"
        "strconv"
)

type Sample struct {
//var Sample struct {
      Time string
      //value []byte
      //time time.Time
      Value map[string]interface{}
      Average map[string]interface{}
      //Value string
}

//type Sample [string][]byte

// Stores the cookie information
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

func printResults(interval int){
	concurrentMap.RLock()
	//name := concurrentMap.cookieMap[formCookie]
	sample := concurrentMap.target
        //avgResult := make([]string, len(sample))
        //dsf := slice[len(sample)]string
        for keys, values := range sample {
          log.Infof("key: %v", keys)
          //log.Infof("value: %v", values)
          for i := 0; i < len(values); i++ {
            //fmt.Printf("value at %v: %v\n", i, values[i])
            //fmt.Println(values[i].Value)
            //TODO print this better????? finsh this up
            fmt.Println(values[i].Average)
            //type amounts struct {
            //      Name  string
            //      Amount string
            //}
            //var jsonResults []amounts
//            var jsonResults map[string]interface{}
//            err := json.Unmarshal([]byte(values[i].Value), &jsonResults)
//            //jsonResult, err := json.Marshal(values[i].Value)
//            if err != nil {
//              log.Errorf("%v", err)
//            }
//            //fmt.Println(jsonResults)
//            for reqName, amount := range jsonResults {
//              //safd
//              fmt.Printf("%s: ", reqName)
//              fmt.Println(amount)
//            }
            //fmt.Printf("total: %s\n", jsonResults["Total"])
            //fmt.Println(jsonResult)
            //err := json.Marshal(values[i].Value, &jsonResults)
            //for reqName, amount := range values[i].Value {
            //  fmt.Println(reqName)
            //  fmt.Println(amount)
            //}
//            reqName := values[i].Value
//            //reqName := Keys(values[i].Value)
//            for i := 0; i < len(reqName); i++ {
//            //for reqName, amount := range values[i].Value {
//              fmt.Println(reqName)
//              //fmt.Println(reqName.Get(i))
//            }
          }
        }
        //fmt.Printf("test %v \n", concurrentMap.target)
        time.Sleep(time.Second)
        //fmt.Printf("test %v \n", concurrentMap.target["http://localhost:9090/monitor"])
        //Print overall average
	concurrentMap.RUnlock()
}

//
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
                        //var jsonResults map[string]interface{}
                        err := json.Unmarshal([]byte(respBody), &jsonResults)
                        if err != nil {
                          log.Errorf("%v", err)
                        }
                        //avgJson := make(map[string]string)
                        //var avgJson map[string]interface{}
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
                            //avg := ((strconv.Atoi(firstValue) + strconv.Atoi(secondValue)) / sampleSec)
                            avg := ((firstValue - secondValue) / sampleSec)
                            //avg := ((int(value) + int(pastJson[key])) / sampleSec)
                            pastJson[key] = avg
                            //avgJson[key] = avg
                            //fmt.Printf("avgKson %v\n", key)
                            //fmt.Printf("avgValu %s\n", value)
                            //fmt.Printf("avgValu %s\n", value)
                          }
                        }
                        //keys = []int{}
                        //for k, value := range  {
                        //      keys = append(keys, k)
                        //}
                        newSample := Sample{fmt.Sprintf("%v", time.Now()), jsonResults, pastJson}
                        //newSample := Sample{fmt.Sprintf("%v", time.Now()), jsonResults, avgJson}
                        //newSample := Sample{fmt.Sprintf("%v", time.Now()), jsonResults, pastJson}
                        sampleSlice = append(sampleSlice, newSample)
                        concurrentMap.target[url] = sampleSlice
                        concurrentMap.Unlock()
		}
                return jsonResults
}

func collectStats(url string, sampleSec int) {
        //firstTime := 
        firstResult := request(url, sampleSec, nil)
	for {
                //TODO compute averages
		time.Sleep(time.Duration(sampleSec) * time.Second)
                firstResult = request(url, sampleSec, firstResult)
		//response, err := http.Get(url)
		//if err != nil {
		//	fmt.Println(err)
		//} else {
		//	//fmt.Println(response)
                //        //log.Infof("Requested url: %s", url)
                //        //fmt.Println(url)
                //        body, _ := ioutil.ReadAll(response.Body)
                //        //body, err := ioutil.ReadAll(response.Body)
	        //        respBody := string(body)
		//	//log.Infof("response body: %s", respBody)
                //        log.Infof("Requested url: %s. Response body: %s", url, respBody)
                //        concurrentMap.Lock()
                //        sampleSlice := concurrentMap.target[url]
                //        //if sampleSlice == nil {
                //        //        //asdf
                //        //}
                //        var jsonResults map[string]interface{}
                //        err := json.Unmarshal([]byte(respBody), &jsonResults)
                //        //jsonResult, err := json.Marshal(values[i].Value)
                //        if err != nil {
                //          log.Errorf("%v", err)
                //        }
                //        newSample := Sample{fmt.Sprintf("%v", time.Now()), jsonResults, currentaverage}
                //        //newSample := Sample{fmt.Sprintf("%v", time.Now()), respBody}
                //        //fmt.Printf("new slice %s\n", newSample)
                //        //fmt.Printf("new slice %s\n", sampleSlice)
                //        //newSample := Sample(time.Now(), respBody)
                //        sampleSlice = append(sampleSlice, newSample)
                //        concurrentMap.target[url] = sampleSlice
                //        //fmt.Printf("sample slice %s\n", sampleSlice)
                //        //sampleSlice = Append(sampleSlice, newSample)
                //        concurrentMap.Unlock()
		//	//fmt.Println(respBody)
		//	//fmt.Println(response.Body)
		//}
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
	//TODO make the targets: split it on the commas
}
