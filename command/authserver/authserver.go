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
	"fmt"
	"html/template"
	"net"
	"os/exec"
	"strings"
	"sync"
	"time"
        "flag"
        "io/ioutil"
        //"ioutil"
	//"json"
        "encoding/json"
	//"filepath"
)

var templatesFolder string
var templatesSlice []string

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

func buildMap(loadfile string) {
  file, fileErr := ioutil.ReadFile(loadfile)
  if fileErr != nil {
    fmt.Println("error:", fileErr)
  }
	concurrentMap.Lock()
        //move this above
        //b, err := json.Marshal(concurrentMap.cookieMap)
         err := json.Unmarshal(file, &concurrentMap.cookieMap)
	concurrentMap.Unlock()
  if err != nil {
    fmt.Println("error:", err)
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
            fmt.Println("error:", err)
        }
        //fmt.Println(b)
        //os.Stdout.Write(b)
        writeError := ioutil.WriteFile("backup/backup.bak", b, 0644)
        if writeError != nil {
          os.Exit(0)
        }
	return
}

func setPath(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	formCookie := r.FormValue("cookie")
	formName := r.FormValue("name")
	if len(formCookie) == 0 || len(formName) == 0 {
		missingCookie := ""
		missingName := ""
		if len(formCookie) == 0 {
			missingCookie = "Cookie is missing"
		}
		if len(formName) == 0 {
			missingName = "Name is missing"
		}
		info := &Information{
			Name:   missingName,
			Cookie: missingCookie,
		}
		malformedRequest(w, r, info)
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
go backupServer(){
}

func main() {
	defer log.Flush()
	dumpfile := flag.String("dumpfile", "backup", "This is the authserver dump file")
	//backupInterval := flag.Int("checkpoint-interval", 10, "This is the authserver backup interval")
	flag.Parse()
	ief, err0 := net.InterfaceByName("eth0")
	if err0 != nil {
		//log.Fatal(err)
	}
	addrs, err1 := ief.Addrs()
	if err1 != nil {
		//log.Fatal(err)
	}
	//fmt.Println("HERE:")
	//fmt.Println(addrs)
	//fmt.Println(addrs[0])
	ipAddr := ""
	if addrs != nil {
		theIP := fmt.Sprintf("%s", addrs[0])
		ipAddr = fmt.Sprintf("%s", strings.Split(theIP, "/")[0])
	} else {
		ipAddr = "localhost"
	}
	fmt.Println(ipAddr)
	fmt.Println(*dumpfile)
        loadingFile := fmt.Sprintf("backup/%s.bak", *dumpfile)
        buildMap(loadingFile)
	http.HandleFunc("/get", getPath)
	http.HandleFunc("/set", setPath)
	http.HandleFunc("/", errorer)
	err := http.ListenAndServe(":9090", nil)
	//err := http.ListenAndServe(portString, nil)
	if err != nil {
		//log.Errorf("Server Failed: %s", err)
		os.Exit(1)
	}
}
