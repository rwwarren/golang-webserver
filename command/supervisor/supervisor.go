// (C) Ryan Warren 2015
// Supervisor
//
// Supervisor monitors servers and makes sure that they
// are still running by checking every checkpoint-interval
// seconds and reloads the server if needs be
//
// There are a couple flags for this program:
// --log is the logger configuration file
// --port-range is the start and end port of servers
// --checkpoint-interval interval between checking services
// --dumpfile is a backup file of the servers that are being
// monitored by the supervisor

package main

import (
	log "../../seelog-master/"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	//"os/signal"
)

// Stores the port information
var concurrentMap struct {
	sync.RWMutex
	portMap []Ports
	//TODO use size for something
	size int
}

//
type Ports struct {
	PortNumber int
	IsUsed     bool
}

//
type configs struct {
	Command     []string
	Output      string
	Error       string
	PID         int
	CurrentPort string
	//TODO PID, PORT
	//TODO output all map to JSON then reload it
}

// Initalizes the concurrent map with the port list
func init() {
	concurrentMap = struct {
		sync.RWMutex
		portMap []Ports
		size    int
	}{portMap: make([]Ports, 9999)}
}

//
func killBackups(killProcesses *configs) {
	if killProcesses.PID != 0 {
		killProcess(killProcesses.PID)
		killProcesses.PID = 0
	}
}

//
func killProcess(pid int) {
	process, _ := os.FindProcess(pid)
	process.Kill()
	process.Wait()
}

//
func checkAlive(pid int) bool {
	if pid == 0 {
		return false
	}
	process, _ := os.FindProcess(pid)
	//TODO get process STAT
	newerr := process.Signal(syscall.Signal(0))
	fmt.Println("newerr")
	fmt.Println(newerr)
	myPid := fmt.Sprintf("%v", pid)
	cmd, err := exec.Command("/bin/ps", "axo pid,stat | grep", myPid).Output()
	if err != nil {
		fmt.Printf("exit status???? %s\n", err)
	}
	fmt.Printf("CMD output %s \n", cmd)
	n := len(cmd)
	command := string(cmd[:n])
	return !strings.Contains(command, "Z")
}

//
func launch(currentConfig *configs, thepath string) {
	size := len(currentConfig.Command) - 1
	args := make([]string, size)
	var foundPort string
	for i := 0; i < size; i++ {
		currentCommand := currentConfig.Command[i+1]
		if strings.Contains(currentCommand, "{{port}}") {
			foundPort = getFreePort()
			log.Infof("Starting on port: %s", foundPort)
			currentCommand = strings.Replace(currentCommand, "{{port}}", foundPort, 1)
		} else if strings.Contains(currentCommand, "--port=") {
			found := strings.SplitAfter(currentCommand, "--port=")
			foundPort = found[1]
			concurrentMap.Lock()
			for i := range concurrentMap.portMap {
				thePort, _ := strconv.Atoi(foundPort)
				concurrentMap.portMap[i].PortNumber = thePort
				concurrentMap.portMap[i].IsUsed = true
				log.Infof("Marking port as used: %s", foundPort)
				break
			}
			concurrentMap.Unlock()
		}
		args[i] = currentCommand
	}
	program := fmt.Sprintf("%s", currentConfig.Command[0])
	cmd := exec.Command(program, args...)
	outfile, outerr := os.Create(currentConfig.Output)
	if outerr != nil {
		panic(outerr)
	}
	defer outfile.Close()
	errfile, errerr := os.Create(currentConfig.Error)
	if errerr != nil {
		panic(errerr)
	}
	defer errfile.Close()
	//TODO send output to a command line output
	cmd.Stdout = outfile
	cmd.Stderr = errfile
	fmt.Println(cmd)
	err := cmd.Start()
	fmt.Printf("ProcessID: %v\n", cmd.Process.Pid)
	if err != nil {
		log.Critical(err)
	}
	currentConfig.PID = cmd.Process.Pid
	currentConfig.CurrentPort = foundPort
}

//
func supervise(currentConfig *configs, thepath string, checkoutInterval int) {
	for {
		alive := checkAlive(currentConfig.PID)
		fmt.Printf("ProcessID Checking: %v\n", currentConfig.PID)
		if alive {
			time.Sleep(time.Duration(checkoutInterval) * time.Second)
		} else {
			launch(currentConfig, thepath)
		}
	}
}

// Gets the next free port for the next server to load
func getFreePort() string {
	concurrentMap.Lock()
	var portNum int
	for i := range concurrentMap.portMap {
		if !concurrentMap.portMap[i].IsUsed {
			concurrentMap.portMap[i].IsUsed = true
			portNum = concurrentMap.portMap[i].PortNumber
			concurrentMap.Unlock()
			log.Infof("portnum %v . port string %s \n", portNum, strconv.Itoa(portNum))
			return strconv.Itoa(portNum)
		}
	}
	concurrentMap.Unlock()
	return strconv.Itoa(portNum)
}

// Builds the list of used ports
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
	for i := 0; i <= total; i++ {
		concurrentMap.portMap[i] = Ports{(min + i), false}
	}
	concurrentMap.size = total
	concurrentMap.Unlock()
}

// Gets the file as a []byte that will be used to load
func getLoadFile(loadingFile string) []byte {
	fileBytes, err := ioutil.ReadFile(loadingFile)
	if err != nil {
		log.Criticalf("Failed: %s", err)
		os.Exit(1)
	}
	return fileBytes
}

// Gets the list of servers to monitor
func getSupervisionList(loadedFile []byte) []configs {
	var configList []configs
	err := json.Unmarshal(loadedFile, &configList)
	if err != nil {
		log.Critical(err)
	}
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
	//
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
	//fmt.Println("supervision list")
	//fmt.Println(supervisionList)

	//TODO remove this?
	filename := os.Args[0]
	filedirectory := filepath.Dir(filename)
	thepath, err := filepath.Abs(filedirectory)
	if err != nil {
		log.Critical(err)
	}
	loadedFile := getLoadFile(dumpfile)
	additionalBackup := getSupervisionList(loadedFile)
	for key, _ := range additionalBackup {
		killBackups(&additionalBackup[key])
	}
	totalSize := (len(additionalBackup) + len(supervisionList))
	newList := make([]configs, len(supervisionList), totalSize)
	copy(newList, supervisionList)
	supervisionList = newList
	supervisionList = append(supervisionList, additionalBackup...)
	for key, _ := range supervisionList {
		go supervise(&supervisionList[key], thepath, checkoutInterval)
	}
	log.Info("Loading the servers")
	//go monitor()
	//TODO channel to keep servers alive
	//signalChan := make(chan os.Signal, 1)
	//signal.Notify(signalChan, os.Interrupt)
	//go func() {
	//	for _ = range signalChan {
	//		fmt.Println("\nReceived shutdown command. Cleaning up...\n")
	//		//cleanup()
	//		os.Exit(0)
	//		//os.Exit(1)
	//	}
	//}()
	//TODO remove wait group??? or only have it run the first time
	//wg.Wait()
	log.Info("Done! Loaded the servers")
	var i int
	//TODO use this to quit
	for i != 1 {
		fmt.Scan(&i)
		fmt.Println("read number", i, "from stdin")
		//cleanup
	}
	os.Exit(0)
}
