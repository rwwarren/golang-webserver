// (C) Ryan Warren 2015
// Authserver
//
// Authserver. This tracks all the user logged in information
// it stores infomation about the user to make sure that another
// server can see if the user is logged in or not
//
// There are a couple flags for this program:
// "-checkpoint-interval" is the authserver backup interval
// "-dumpfile" is the authserver backup file
// "-log" is the logger configuration file
// "-port" is the auth server port

package main

import (
	log "../../seelog-master/"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
)

var templatesFolder string
var templatesSlice []string
var done bool
var loadingFile string
var backupFile string

// Stores the cookie information
var concurrentMap struct {
	sync.RWMutex
	cookieMap map[string]string
}

// Initalizes the authserver with the important user storage things
func init() {
	templatesFolder = "templates"
	templatesSlice = append(templatesSlice, fmt.Sprintf("%s/template.html", templatesFolder))
	concurrentMap = struct {
		sync.RWMutex
		cookieMap map[string]string
	}{cookieMap: make(map[string]string)}
}

// Information about the user, their cookie and username
type Information struct {
	Name   string
	Cookie string
}

// Builds the user map based off the file passed in.
// It must be in json format. Will not crash if the
// file does not exist
func buildMap(loadFile string) {
	file, fileErr := ioutil.ReadFile(loadFile)
	if fileErr != nil {
		log.Errorf("file error: %s", fileErr)
		return
	}
	concurrentMap.Lock()
	err := json.Unmarshal(file, &concurrentMap.cookieMap)
	concurrentMap.Unlock()
	if err != nil {
		log.Errorf("build error: %s", err)
		os.Exit(1)
	}
}

// If the request is not formed correctly, this will return a 400 error to the user
func malformedRequest(w http.ResponseWriter, r *http.Request, missingInfo *Information) {
	w.WriteHeader(400)
	malformedPageTemplatesSlice := make([]string, len(templatesSlice))
	copy(malformedPageTemplatesSlice, templatesSlice)
	malformedPageTemplatesSlice = append(malformedPageTemplatesSlice, fmt.Sprintf("%s/malformed.html", templatesFolder))
	var malformedPage = template.Must(template.New("MalformedPage").ParseFiles(malformedPageTemplatesSlice...))
	malformedPage.ExecuteTemplate(w, "template", missingInfo)
	return
}

// Returns the name of the user based on the cookie that is passed in
func getPath(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	formCookie := r.FormValue("cookie")
	if len(formCookie) == 0 {
		missingCookie := ""
		missingName := "Name is missing"
		info := &Information{
			Name:   missingName,
			Cookie: missingCookie,
		}
		malformedRequest(w, r, info)
		return
	}
	concurrentMap.RLock()
	name := concurrentMap.cookieMap[formCookie]
	concurrentMap.RUnlock()
	getPageTemplatesSlice := make([]string, len(templatesSlice))
	copy(getPageTemplatesSlice, templatesSlice)
	getPageTemplatesSlice = append(getPageTemplatesSlice, fmt.Sprintf("%s/get.html", templatesFolder))
	var getPage = template.Must(template.New("GetPage").ParseFiles(getPageTemplatesSlice...))
	getPage.ExecuteTemplate(w, "template", name)
	return
}

// Sets the name of the user based on the cookie and the name passed in
func setPath(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	formCookie := r.FormValue("cookie")
	formName := r.FormValue("name")
	if len(formCookie) == 0 {
		missingCookie := ""
		missingName := ""
		if len(formCookie) == 0 {
			missingCookie = "Cookie is missing"
		}
		info := &Information{
			Name:   missingName,
			Cookie: missingCookie,
		}
		malformedRequest(w, r, info)
		return
	} else if len(formName) == 0 {
		logout(formCookie)
		return
	}
	setPageTemplatesSlice := make([]string, len(templatesSlice))
	copy(setPageTemplatesSlice, templatesSlice)
	setPageTemplatesSlice = append(setPageTemplatesSlice, fmt.Sprintf("%s/set.html", templatesFolder))
	var setPage = template.Must(template.New("SetPage").ParseFiles(setPageTemplatesSlice...))
	setPage.ExecuteTemplate(w, "template", "")
	concurrentMap.Lock()
	concurrentMap.cookieMap[formCookie] = formName
	concurrentMap.Unlock()
	return
}

// Returns an error 404 that the page is not found
func errorer(w http.ResponseWriter, r *http.Request) {
	log.Infof("Error, url not found for authserver: %s", r.URL)
	w.WriteHeader(404)
	errorTemplatesSlice := make([]string, len(templatesSlice))
	copy(errorTemplatesSlice, templatesSlice)
	errorTemplatesSlice = append(errorTemplatesSlice, fmt.Sprintf("%s/404.html", templatesFolder))
	var errorPage = template.Must(template.New("ErrorPage").ParseFiles(errorTemplatesSlice...))
	errorPage.ExecuteTemplate(w, "template", "")
	return
}

// Backs up the server into json format
func backupServer(backupInterval int) {
	for !done {
		time.Sleep(time.Duration(backupInterval) * time.Second)
		if _, fileErr := os.Stat(loadingFile); fileErr == nil && !strings.Contains(loadingFile, ".bak") {
			loadingFile = fmt.Sprintf("%s.bak", loadingFile)
			deleteBackup(loadingFile)
			writeBackup(loadingFile)
			buildMap(loadingFile)
		} else if strings.Contains(loadingFile, ".bak") {
			deleteBackup(loadingFile)
			writeBackup(loadingFile)
			buildMap(loadingFile)
		}
	}
}

// Cleans up the server and saves it before quitting, after the interrupt command
// is recieved
func cleanup() {
	writeBackup(backupFile)
	deleteBackup(loadingFile)
	log.Info("Cleanup Complete")
	fmt.Println("Cleanup Complete")
}

// Writes the user information to the file passed in
// to make sure in case of shutdown there is a copy
func writeBackup(BackupFilename string) {
	concurrentMap.RLock()
	backup := make(map[string]string)
	for k, v := range concurrentMap.cookieMap {
		backup[k] = v
	}
	concurrentMap.RUnlock()
	b, err := json.Marshal(backup)
	if err != nil {
		log.Errorf("error reading json: %s", err)
		os.Exit(1)
	}
	writeError := ioutil.WriteFile(BackupFilename, b, 0644)
	if writeError != nil {
		log.Errorf("error: %s", writeError)
		os.Exit(1)
	}
	log.Infof("Backup complete to: %s", BackupFilename)
}

// Deletes the copy of the backup file
func deleteBackup(filename string) {
	if strings.Contains(filename, ".bak") {
		log.Infof("Deleting backup file: %s", filename)
		os.Remove(filename)
	}
}

// Logs the user out of the server
func logout(uuid string) {
	concurrentMap.Lock()
	delete(concurrentMap.cookieMap, uuid)
	concurrentMap.Unlock()
	log.Infof("Logging user out: %s", uuid)
	return
}

// Main function of the authserver. Runs by default on port 9090
func main() {
	defer log.Flush()
	port := flag.Int("port", 9090, "Set the server port, default port: 9090")
	dumpfile := flag.String("dumpfile", "backup", "This is the authserver dump file")
	backupInterval := flag.Int("checkpoint-interval", 10, "This is the authserver backup interval")
	logFile := flag.String("log", "logConfig", "This is the logger configuration file")
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
	ief, err0 := net.InterfaceByName("eth0")
	if err0 != nil {
		log.Critical(err0)
	}
	addrs, err1 := ief.Addrs()
	if err1 != nil {
		log.Critical(err1)
	}
	ipAddr := ""
	if addrs != nil {
		theIP := fmt.Sprintf("%s", addrs[0])
		ipAddr = fmt.Sprintf("%s", strings.Split(theIP, "/")[0])
	} else {
		ipAddr = "localhost"
	}
	var portString = fmt.Sprintf(":%d", *port)
	log.Infof("IpAddress and port: %s%s", ipAddr, portString)
	loadingFile = fmt.Sprintf("backup/%s", *dumpfile)
	backupFile = fmt.Sprintf("backup/%s", *dumpfile)
	buildMap(loadingFile)
	done = false
	go backupServer(*backupInterval)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for _ = range signalChan {
			fmt.Println("\nReceived shutdown command. Cleaning up...\n")
			cleanup()
			os.Exit(1)
		}
	}()
	http.HandleFunc("/get", getPath)
	http.HandleFunc("/set", setPath)
	http.HandleFunc("/", errorer)
	err := http.ListenAndServe(portString, nil)
	if err != nil {
		log.Errorf("Server Failed: %s", err)
		os.Exit(1)
	}
	done = true
}
