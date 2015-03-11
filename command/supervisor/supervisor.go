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
	"os/exec"
	"strings"
	//test
//	"bytes"
	//
	"encoding/json"
	"flag"
	"fmt"
	//"io/ioutil"
	//"net/http"
	"os"
	"strconv"
	//"strings"
	"sync"
	//"time"
	"io/ioutil"
	"path/filepath"
        //"runtime"
)

// Stores the port information
var concurrentMap struct {
	sync.RWMutex
	portMap []Ports
	size    int
	//portMap map[int]bool
}

type Ports struct {
	PortNumber int
	IsUsed     bool
}

type configs struct {
	Command []string
	Output  string
	Error   string
        PID     int
        CurrentPort string
        //TODO PID, PORT
        //TODO output all map to JSON then reload it
}

// Initalizes the
func init() {
	concurrentMap = struct {
		sync.RWMutex
		portMap []Ports
		size    int
		//portMap map[int]bool
		//}{}
		//}{portMap: make(map[int]bool)}
	}{portMap: make([]Ports, 9999)}
}

//
func supervise(currentConfig configs, thepath string, wg *sync.WaitGroup) {
	//supervise the server
	size := len(currentConfig.Command) - 1
	args := make([]string, size)
	//fmt.Println(size)
        var foundPort string
	for i := 0; i < size; i++ {
		//currentCommand := currentConfig.Command[i]
		currentCommand := currentConfig.Command[i + 1]
		if strings.Contains(currentCommand, "{{port}}") {
			foundPort = getFreePort()
                        fmt.Printf("FOUND PORT: %s\n", foundPort)
			currentCommand = strings.Replace(currentCommand, "{{port}}", foundPort, 1)
		}
		fmt.Printf("at this sport: %s\n", currentCommand)
		args[i] = currentCommand
	}
        //attr := new(os.ProcAttr)
        //cmd, err := os.StartProcess(currentConfig.Command[0], args, attr)
        //if err != nil {
        //  log.Critical(err)
        //}
        //fmt.Printf("cmd!!! %v \n", cmd)
        program := fmt.Sprintf("%s", currentConfig.Command[0])
        //program := fmt.Sprintf("%s/%s", thepath, currentConfig.Command[0])
	cmd := exec.Command(program, args...)
        outfile, outerr := os.Create(currentConfig.Output)
        //outfile, outerr := os.Create("./out.txt")
        if outerr != nil {
            panic(outerr)
        }
        defer outfile.Close()
        errfile, errerr := os.Create(currentConfig.Error)
        //outfile, outerr := os.Create("./out.txt")
        if errerr != nil {
            panic(outerr)
        }
        defer errfile.Close()
        //TODO send output to a command line output
//        cmd.Stdout = os.Stdout
        cmd.Stdout = outfile
        cmd.Stderr = errfile
	fmt.Println(cmd)
	//cmd := exec.Command(currentConfig.Command[0], currentConfig.Command[1:size]...)
	//testing
	//var out bytes.Buffer
	//cmd.Stdout = &out
	//
	err := cmd.Start()
	//err := cmd.Run()
	fmt.Printf("ProcessID: %v\n", cmd.Process.Pid)
	//fmt.Println(cmd.Process.Pid)
	if err != nil {
		log.Critical(err)
	}
        concurrentMap.Lock()
        currentConfig.PID = cmd.Process.Pid
        currentConfig.CurrentPort = foundPort
        concurrentMap.Unlock()
        wg.Done()
	//fmt.Printf("in all caps: %q\n", out.String())
        //output <- out.String()
//        for {
//        }
        //cmd.Wait()

}

func getFreePort() string {
	concurrentMap.Lock()
	var portNum int
	for i := range concurrentMap.portMap {
		if !concurrentMap.portMap[i].IsUsed {
			//fmt.Printf("THIS LOCATION %v \n", i)
			concurrentMap.portMap[i].IsUsed = true
			portNum = concurrentMap.portMap[i].PortNumber
	                concurrentMap.Unlock()
                        fmt.Printf("portnum %v . port string %s \n", portNum, strconv.Itoa(portNum))
                        //return string(portNum)
                        return strconv.Itoa(portNum)
		}
	}
	concurrentMap.Unlock()
	return strconv.Itoa(portNum)
	//return string(portNum)
}

func buildPorts(ports []string) {
	min, minerr := strconv.Atoi(ports[0])
	max, maxerr := strconv.Atoi(ports[1])
	if minerr != nil || maxerr != nil {
		// handle error
		fmt.Println(minerr)
		fmt.Println(maxerr)
		os.Exit(2)
	}
	concurrentMap.Lock()
	total := (max - min)
	//fmt.Println(total)
	//myMap := make([]Ports, total)
	//concurrentMap.portMap = myMap
	for i := 0; i <= total; i++ {
		//for currentPort := min; currentPort <= max; currentPort++ {
		//myMap[i] = Ports{(min + i), false}
		concurrentMap.portMap[i] = Ports{(min + i), false}
		//fmt.Println(concurrentMap.portMap[i])
		//concurrentMap.portMap[currentPort] = false
	}
	concurrentMap.size = total
	concurrentMap.Unlock()

}

func getLoadFile(loadingFile string) []byte {
	fileBytes, err := ioutil.ReadFile(loadingFile)
	if err != nil {
		log.Criticalf("Failed: %s", err)
		os.Exit(1)
	}
	return fileBytes
}

func getSupervisionList(loadedFile []byte) []configs {
	var configList []configs
	err := json.Unmarshal(loadedFile, &configList)
	if err != nil {
		log.Critical(err)
	}
	//for key, val := range configList {
	//	fmt.Println(key)
	//	fmt.Println(val)
	//}
	return configList
}

// Main function of the supervisor
func main() {
	defer log.Flush()
	logFile := flag.String("log", "logConfig", "This is the logger configuration file")
	portRange := flag.String("port-range", "8080-9090", "This is the port range")
	dumpLoc := flag.String("dumpfile", "backup.bak", "This is the dumpfile")
        //TODO below change to pipe in from command line
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
	log.Infof("dumpfile Flag: %s", dumpfile)
	log.Infof("checkpoint interval Flag: %v", checkoutInterval)
	log.Infof("loading file Flag: %s", loadingFile)
	buildPorts(portsList)
	loadedString := getLoadFile(loadingFile)
	supervisionList := getSupervisionList(loadedString)
	//fmt.Println(supervisionList)

        //TODO remove this?
	filename := os.Args[0]
	filedirectory := filepath.Dir(filename)
	thepath, err := filepath.Abs(filedirectory)
	if err != nil {
		log.Critical(err)
	}
	fmt.Println(thepath)
        //

        wg := new(sync.WaitGroup)
        amount := len(supervisionList)
        wg.Add(amount)
        //ch := make(chan string)
	for _, val := range supervisionList {
		//for key, val := range supervisionList {
		//fmt.Println(key)
		fmt.Println(val)
		//go supervise(val, ch)
		//supervise(val, ch)
		go supervise(val, thepath, wg)
                //go test("/bin/ls")
                //test("/bin/ls")
//		supervise(val, thepath)
		//supervise(val)
	}
        //for {
        //  //fmt.Print(runtime.NumGoroutine())

        //}
        //for i := range ch {
        //     fmt.Printf("%v ", i)
        //}

        fmt.Println("Loading the servers")
        wg.Wait()
        fmt.Println("Done! Loaded the servers")
	//strings.Replace on {{port}}
        //time.Sleep(5 * time.Second)
        fmt.Println("sleep done")
        //for {
        //}
}
