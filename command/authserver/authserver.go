// (C) Ryan Warren 2015
// Authserver
//
//

package main

import (
	"net/http"
	"os"
	//"os/exec"
	log "../../seelog-master/"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"net"
	"os/exec"
	"strings"
	"sync"
	"time"
	//"ioutil"
	//"json"
	"encoding/json"
	//"filepath"
)

var templatesFolder string
var templatesSlice []string
var done bool
var loadingFile string

// Stores the cookie information
var concurrentMap struct {
	sync.RWMutex
	cookieMap map[string]string
}

func init() {
	templatesFolder = "templates"
	templatesSlice = append(templatesSlice, fmt.Sprintf("%s/template.html", templatesFolder))
	concurrentMap = struct {
		sync.RWMutex
		cookieMap map[string]string
	}{cookieMap: make(map[string]string)}
}

type Information struct {
	Name   string
	Cookie string
}

func buildMap() {
//func buildMap(loadfile string) {
	file, fileErr := ioutil.ReadFile(loadingFile)
	if fileErr != nil {
		log.Errorf("file error: %s", fileErr)
		//os.Exit(1)
                return
	}
	concurrentMap.Lock()
	//move this above
	//b, err := json.Marshal(concurrentMap.cookieMap)
	err := json.Unmarshal(file, &concurrentMap.cookieMap)
	concurrentMap.Unlock()
	if err != nil {
		log.Errorf("build error: %s", err)
		os.Exit(1)
	}
	//if the read-back is successful, the backup file should be deleted.
	//fmt.Printf("%+v", animals)
}

// Set and returns the cookie from the request
func SetCookie(w http.ResponseWriter, r *http.Request) *http.Cookie {
	checkCookie, cookieError := r.Cookie("uuid")
	if cookieError == nil {
		log.Infof("Cookie is already set: %s", checkCookie.Value)
		return checkCookie
	}
	uuid, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Infof("Error something went wrong with uuidgen: %s", err)
		os.Exit(1)
	}
	uuidLen := len(uuid) - 1
	uuidString := string(uuid[:uuidLen])
	log.Infof("Setting cookie with UUID: %s", uuidString)
	cookie := &http.Cookie{Name: "uuid", Value: uuidString, Expires: time.Now().Add(356 * 24 * time.Hour), HttpOnly: true}
	http.SetCookie(w, cookie)
	return cookie
}

func malformedRequest(w http.ResponseWriter, r *http.Request, missingInfo *Information) {
	w.WriteHeader(400)
	malformedPageTemplatesSlice := make([]string, len(templatesSlice))
	copy(malformedPageTemplatesSlice, templatesSlice)
	malformedPageTemplatesSlice = append(malformedPageTemplatesSlice, fmt.Sprintf("%s/malformed.html", templatesFolder))
	var malformedPage = template.Must(template.New("MalformedPage").ParseFiles(malformedPageTemplatesSlice...))
	malformedPage.ExecuteTemplate(w, "template", missingInfo)
	return
}

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
	//
	getPageTemplatesSlice := make([]string, len(templatesSlice))
	copy(getPageTemplatesSlice, templatesSlice)
	getPageTemplatesSlice = append(getPageTemplatesSlice, fmt.Sprintf("%s/get.html", templatesFolder))
	var getPage = template.Must(template.New("GetPage").ParseFiles(getPageTemplatesSlice...))
	getPage.ExecuteTemplate(w, "template", name)
	return
}

func setPath(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	formCookie := r.FormValue("cookie")
	formName := r.FormValue("name")
	if len(formCookie) == 0 {
	//if len(formCookie) == 0 || len(formName) == 0 {
		missingCookie := ""
		missingName := ""
		if len(formCookie) == 0 {
			missingCookie = "Cookie is missing"
		}
		//if len(formName) == 0 {
		//	missingName = "Name is missing"
		//}
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
	//printRequests(r)
	//log.Info("Error, url not found: These are not the URLs you are looking for.")
	//w.WriteHeader(404)
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

func errorer(w http.ResponseWriter, r *http.Request) {
	//printRequests(r)
	//log.Info("Error, url not found: These are not the URLs you are looking for.")
	w.WriteHeader(404)
	errorTemplatesSlice := make([]string, len(templatesSlice))
	copy(errorTemplatesSlice, templatesSlice)
	errorTemplatesSlice = append(errorTemplatesSlice, fmt.Sprintf("%s/404.html", templatesFolder))
	var errorPage = template.Must(template.New("ErrorPage").ParseFiles(errorTemplatesSlice...))
	errorPage.ExecuteTemplate(w, "template", "")
	return
}

//called go routine
//https://gobyexample.com/goroutines
func backupServer(backupInterval int) {
//func backupServer(backupInterval int, loadingFile string) {
	//func backupServer(done chan bool){
	for !done {
		time.Sleep(time.Duration(backupInterval) * time.Second)
		writeBackup()
		buildMap()
		//deleteBackup()
	}
}

func writeBackup() {
//func writeBackup(loadingFile string) {
	//TODO fix this
	concurrentMap.RLock()
	backup := make(map[string]string)
	//backup := make(map[string]string, len(concurrentMap.cookieMap))
	//copy(backup, concurrentMap.cookieMap)
	for k, v := range concurrentMap.cookieMap {
		backup[k] = v
	}
	//move this above
	concurrentMap.RUnlock()
	b, err := json.Marshal(backup)
	if err != nil {
		log.Errorf("error: %s", err)
		os.Exit(1)
	}
      //TODO fix this!
	if _, fileErr := os.Stat(loadingFile); fileErr == nil && !strings.Contains(loadingFile, ".bak") {
              log.Info("got here")
	      loadingFile = fmt.Sprintf("%s.bak", loadingFile)
              log.Info("got here too")
              deleteBackup()
	//} else {
        //  log.Info(fileErr)
        }
	writeError := ioutil.WriteFile(loadingFile, b, 0644)
	if writeError != nil {
		log.Errorf("error: %s", writeError)
		os.Exit(1)
	}
}

func deleteBackup() {
//func deleteBackup(loadingFile string) {
	os.Remove(loadingFile)
        log.Info("here")
}

func logout(uuid string) {
        log.Infof("Here is the uuid: %s", uuid)
	concurrentMap.Lock()
	//concurrentMap.cookieMap[formCookie] = formName
	delete(concurrentMap.cookieMap, uuid)
	concurrentMap.Unlock()
        log.Infof("Logging user out: %s", uuid)
        return
}

func main() {
	defer log.Flush()
	port := flag.Int("port", 9090, "Set the server port, default port: 9090")
	dumpfile := flag.String("dumpfile", "backup", "This is the authserver dump file")
	backupInterval := flag.Int("checkpoint-interval", 1, "This is the authserver backup interval")
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
	buildMap()
	done = false
	go backupServer(*backupInterval)
	//go backupServer(*backupInterval, loadingFile)
	//go backupServer(done)
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
