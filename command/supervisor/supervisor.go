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
	"strconv"
	"strings"
	"sync"
	"time"
	"os/signal"
        "io"
        "bytes"
)

// Stores the port supervised server information
var concurrentMap struct {
	sync.RWMutex
	portMap []Ports
        configList []configs
}

// Port information for the supervisor
type Ports struct {
	PortNumber int
	IsUsed     bool
}

// Server configuration struct
type configs struct {
	Command     []string
	Output      string
	Error       string
	PID         int
	CurrentPort string
}

// Initalizes the concurrent map with the port list
func init() {
	concurrentMap = struct {
		sync.RWMutex
		portMap []Ports
                configList []configs
	}{}
}

// Helps kill processes based on the pid
func killBackups(killProcesses *configs) {
	if killProcesses.PID != 0 {
		killProcess(killProcesses.PID)
		killProcesses.PID = 0
	}
}

// Kills he process based on the pid
func killProcess(pid int) {
	process, _ := os.FindProcess(pid)
	process.Kill()
	process.Wait()
}

// Checks if the process is alive
func checkAlive(pid int) bool {
	if pid == 0 {
		return false
	}
      c1 := exec.Command("/bin/ps", "axo pid,stat")
    myPid := fmt.Sprintf("%v", pid)
    c2 := exec.Command("grep", myPid)
    r, w := io.Pipe()
    c1.Stdout = w
    c2.Stdin = r
    var cmd bytes.Buffer
    c2.Stdout = &cmd
    c1.Start()
    c2.Start()
    c1.Wait()
    w.Close()
    c2.Wait()
        command := cmd.String()
	return !strings.Contains(command, "Z")
}

// Launches the service and stores important information
func launch(currentConfig *configs) {
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
		log.Critical(outerr)
	}
	defer outfile.Close()
	errfile, errerr := os.Create(currentConfig.Error)
	if errerr != nil {
		log.Critical(errerr)
	}
	defer errfile.Close()
	cmd.Stdout = outfile
	cmd.Stderr = errfile
	log.Infof("Command exec info: %v", cmd)
	err := cmd.Start()
	log.Infof("ProcessID: %v", cmd.Process.Pid)
	if err != nil {
		log.Critical(err)
	}
        concurrentMap.Lock()
	currentConfig.PID = cmd.Process.Pid
	currentConfig.CurrentPort = foundPort
        concurrentMap.Unlock()
}

// Watches over the server, checking it is up every checkpoint interval
func supervise(currentConfig *configs, checkoutInterval int) {
	for {
		alive := checkAlive(currentConfig.PID)
		log.Infof("ProcessID Checking: %v", currentConfig.PID)
		if alive {
			time.Sleep(time.Duration(checkoutInterval) * time.Second)
		} else {
			launch(currentConfig)
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
        concurrentMap.portMap = make([]Ports, total)
	for i := 0; i < total; i++ {
		concurrentMap.portMap[i] = Ports{(min + i), false}
	}
	concurrentMap.Unlock()
}

// Gets the file as a []byte that will be used to load
func getLoadFile(loadingFile string) []byte {
	fileBytes, err := ioutil.ReadFile(loadingFile)
	if err != nil {
		log.Criticalf("Failed: %s", err)
                return nil
		//os.Exit(1)
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

// Writes the config list to a file
func writeBackup(backupFile string) {
    concurrentMap.RLock()
    backup := &concurrentMap.configList
    concurrentMap.RUnlock()
    b, err := json.Marshal(backup)
    if err != nil {
        log.Critical(err)
        return
    }
    file, err := os.Create(backupFile)
    if err != nil {
        log.Critical(err)
        return
    }
    defer file.Close()
    written, werr := file.Write(b)
    if werr != nil {
        log.Critical(werr)
        return
    }
    log.Infof("Wrote to file: %v", written)
}

// Writes the backup every checkpoint interval
func monitor(backupFile string, checkpointInterval int){
  for {
    time.Sleep(time.Duration(checkpointInterval) * time.Second)
    writeBackup(backupFile)
    log.Info("To quit: Press Q + ENTER")
    fmt.Println("To quit: Press Q + ENTER")
    log.Info("To kill all process group: Press COMMAND + C")
    fmt.Println("To kill all process group: Press COMMAND + C")
  }
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
        concurrentMap.Lock()
        concurrentMap.configList = make([]configs, totalSize)
        concurrentMap.configList = supervisionList
	for key, _ := range concurrentMap.configList {
		go supervise(&concurrentMap.configList[key], checkoutInterval)
	}
        concurrentMap.Unlock()
	log.Info("Loading the servers")
	go monitor(dumpfile, checkoutInterval)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for _ = range signalChan {
			fmt.Println("\nReceived shutdown command. Cleaning up...\n")
                        writeBackup(dumpfile)
			os.Exit(0)
		}
	}()
	log.Info("Done! Loaded the servers")
	var input string
	for input != "Q" {
		fmt.Scan(&input)
		fmt.Println("Qutting based on command from stdin: ", input)
                writeBackup(dumpfile)
	}
	os.Exit(0)
}
