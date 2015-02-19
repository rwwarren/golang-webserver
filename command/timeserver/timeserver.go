// (C) Ryan Warren 2015
// Timeserver
//
// This is a timeserver that shows the time for the current timezone
// Hours, Minutes, and Seconds on the "/time" url.
// All other urls will return with a 404.
//
// Building and running with go can be done with "go build"
// then "./timeserver". This will show up on port 8080 by default.
// If the port is in use there will be an error.
//
// There are command line options.
// "--port PORT_NUMBER" will change the port.
// "-V" will show the version and then quit.
// "-log FILENAME" pass in the log configuration file and it
// will load it from etc/FILENAME.xml

package main

import (
	log "../../seelog-master/"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
        "os/exec"
	"sync"
	"time"
        "math/rand"
        "io/ioutil"
        "strings"
)

type PageInformation struct {
	Name        string
	CurrentTime string
	UTCtime     string
}

var templatesFolder string
var templatesSlice []string
var server string
var maxInbound int
var authTimeout int
var avgResponse int
var deviation int
var logFile string
var version bool
var port int
var templatesFlag string
var inboundRequests struct {
	sync.RWMutex
	currentRequests int
}

// Intitalizes the timeserver with authserver information
func init() {
	authPort := flag.Int("authport", 9090, "This is the authserver default port")
	authHost := flag.String("authhost", "http://localhost", "This is the authserver default host")
	inflight := flag.Int("max-inflight", 0, "Max number of inflight requests")
	port = *flag.Int("port", 8080, "Set the server port, default port: 8080")
	version = *flag.Bool("V", false, "Shows the version of the timeserver")
	logFile = *flag.String("log", "logConfig", "This is the logger configuration file")
	templatesFlag = *flag.String("templates", "templates", "This is the templates folder name")
	authTimeout = *flag.Int("authtimeout-ms", 1000, "This is the authserver timeout")
	avgResponse = *flag.Int("avg-response-ms", 1000, "This is the timeserver avg response time")
	deviation = *flag.Int("deviation-ms", 10, "This is the timeserver deviation")
        flag.Parse()
        server = fmt.Sprintf("%s:%d", *authHost, *authPort)
        inboundRequests = struct {
            sync.RWMutex
            currentRequests int
        }{currentRequests: 0}
        maxInbound = *inflight
}

// Decreases the amount of current inbound requests on the server
func removeInboundRequest(){
  inboundRequests.Lock()
  inboundRequests.currentRequests--
  currentRequests := inboundRequests.currentRequests
  inboundRequests.Unlock()
  log.Infof("Current inbound requests: %d", currentRequests)
}

// Increases the amount of current inbound requests on the server
func addInboundRequest(){
  inboundRequests.Lock()
  inboundRequests.currentRequests++
  currentRequests := inboundRequests.currentRequests
  inboundRequests.Unlock()
  log.Infof("Current inbound requests: %d", currentRequests)
}

// Returns true if you can have more inbound requests
func canHaveMoreInboundRequests() bool {
  inboundRequests.RLock()
  currentInflight := inboundRequests.currentRequests
  inboundRequests.RUnlock()
  return (maxInbound == 0 || currentInflight < maxInbound)
}

// Error page for there being too many inbound requests on the server
func maxInboundError(w http.ResponseWriter, r *http.Request){
	printRequests(r)
	log.Info("Error, too many inbound requests")
	w.WriteHeader(503)
	errorTemplatesSlice := make([]string, len(templatesSlice))
	copy(errorTemplatesSlice, templatesSlice)
	errorTemplatesSlice = append(errorTemplatesSlice, fmt.Sprintf("%s/503.html", templatesFolder))
	var errorPage = template.Must(template.New("ErrorPage").ParseFiles(errorTemplatesSlice...))
	errorPage.ExecuteTemplate(w, "template", "")
        removeInboundRequest()
	return
}

// Deletes the cookie associated with a resopnse and sets a new one
func deleteCookie(w http.ResponseWriter, r *http.Request){
	deletingCookie := &http.Cookie{Name: "uuid", Value: "s", Expires: time.Unix(1, 0), HttpOnly: true}
	http.SetCookie(w, deletingCookie)
}

// Set and returns the cookie from the request
func setCookie(w http.ResponseWriter, r *http.Request) *http.Cookie {
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

// Gets the username associated with the cookie
func getName(uuid string, w http.ResponseWriter, r *http.Request) string {
  getUrl := fmt.Sprintf("%s/get?cookie=%s", server, uuid)
  timeout := time.Duration(time.Duration(authTimeout) * time.Millisecond)
  client := http.Client{
        Timeout: timeout,
  }
  resp, err := client.Get(getUrl)
  if err != nil {
    log.Criticalf("Error getting authserver: %s" , err)
    deleteCookie(w, r)
    return ""
  }
  log.Infof("Response from the authserver: %s", resp)
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  respBody := string(body)
  firstBody := strings.Split(respBody, "<body>")
  firstBodyHalf := firstBody[1]
  secondBody := strings.Split(firstBodyHalf, "</body>")
  secondBodyHalf := secondBody[0]
  finalBody := strings.Trim(secondBodyHalf, "\n ")
  return finalBody
}

// Sends the name associated with the cookie to the authserver
func setName(uuid string, name string){
  log.Infof("setting name with uuid: %s and name: ", uuid, name)
  setUrl := fmt.Sprintf("%s/set?cookie=%s&name=%s", server, uuid, name)
  resp, err := http.Get(setUrl)
  if err != nil {
    log.Criticalf("Error getting authserver: %s" , err)
    os.Exit(1)
  }
  log.Info("Response: %s", resp)

}

// Handles the timeserver which shows the current time
// for the local timezone
func timeHandler(w http.ResponseWriter, r *http.Request) {
        waitTime := rand.NormFloat64() * float64(deviation) + float64(avgResponse)
        if waitTime < 0 {
          waitTime *= -1
        }
        log.Infof("Artifically creating delay for: %v milliseconds", waitTime)
        time.Sleep(time.Duration(waitTime) * time.Millisecond)
        addInboundRequest()
        if !canHaveMoreInboundRequests() {
          maxInboundError(w, r)
          return
        }
	printRequests(r)
	const layout = "3:04:05 PM"
	const UTClayout = "15:04:05 MST"
	personalString := ""
        isLoggedIn, name := checkLogin(w, r)
        if isLoggedIn {
		personalString = fmt.Sprintf(", %s", name)
		log.Debugf("User is logged in as: %s", name)
	}
	timeTemplatesSlice := make([]string, len(templatesSlice))
	copy(timeTemplatesSlice, templatesSlice)
	timeTemplatesSlice = append(timeTemplatesSlice, fmt.Sprintf("%s/time.html", templatesFolder))
	var timeTmpl = template.Must(template.New("time").ParseFiles(timeTemplatesSlice...))
	currentTime := time.Now().Local().Format(layout)
	UTCTime := time.Now().UTC().Format(UTClayout)
	data := &PageInformation{
		Name:        personalString,
		CurrentTime: currentTime,
		UTCtime:     UTCTime,
	}
        removeInboundRequest()
	timeTmpl.ExecuteTemplate(w, "template", data)
	return
}

// Handles errors for when the page is not found
func errorer(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		indexPage(w, r)
		return
	}
	printRequests(r)
	log.Info("Error, url not found: These are not the URLs you are looking for.")
	w.WriteHeader(404)
	errorTemplatesSlice := make([]string, len(templatesSlice))
	copy(errorTemplatesSlice, templatesSlice)
	errorTemplatesSlice = append(errorTemplatesSlice, fmt.Sprintf("%s/404.html", templatesFolder))
	var errorPage = template.Must(template.New("ErrorPage").ParseFiles(errorTemplatesSlice...))
	errorPage.ExecuteTemplate(w, "template", "")
	return
}

// Checks the index page and will render the
// index based on the user being loggedin or not
func indexPage(w http.ResponseWriter, r *http.Request) {
	printRequests(r)
	isLoggedIn, name := checkLogin(w, r)
	if isLoggedIn {
		log.Debug("User is loggedin, going to loggedin page")
		renderIndex(w, name)
		return
	} else {
		log.Debug("User is not loggedin, going to login page")
		renderLogin(w, r)
		return
	}
}

// Renders the page for a loggedin user
func renderIndex(w http.ResponseWriter, name string) {
	indexTemplatesSlice := make([]string, len(templatesSlice))
	copy(indexTemplatesSlice, templatesSlice)
	indexTemplatesSlice = append(indexTemplatesSlice, fmt.Sprintf("%s/index.html", templatesFolder))
	var indexPage = template.Must(template.New("IndexPage").ParseFiles(indexTemplatesSlice...))
	person := &PageInformation{
		Name: name,
	}
	indexPage.ExecuteTemplate(w, "template", person)
}

// Renders the page if there is no name passed into
// the login page
func renderNoNamePage(w http.ResponseWriter) {
	noNameTemplatesSlice := make([]string, len(templatesSlice))
	copy(noNameTemplatesSlice, templatesSlice)
	noNameTemplatesSlice = append(noNameTemplatesSlice, fmt.Sprintf("%s/noNamePage.html", templatesFolder))
	var noNamePage = template.Must(template.New("NoNamePage").ParseFiles(noNameTemplatesSlice...))
	noNamePage.ExecuteTemplate(w, "template", "")
}

// Renders the login page to the website
func renderLogin(w http.ResponseWriter, r *http.Request) {
	loginTemplatesSlice := make([]string, len(templatesSlice))
	copy(loginTemplatesSlice, templatesSlice)
	loginTemplatesSlice = append(loginTemplatesSlice, fmt.Sprintf("%s/loginPage.html", templatesFolder))
	var loginPage = template.Must(template.New("LoginPage").ParseFiles(loginTemplatesSlice...))
	loginPage.ExecuteTemplate(w, "template", "")
}

// Checks the if the user is logged in and if there is a user
// associated with the cookie
func checkLogin(w http.ResponseWriter, r *http.Request) (bool, string) {
	cookie := setCookie(w, r)
	//cookie := CookieManagement.SetCookie(w, r)
	//concurrentMap.RLock()
	//name := concurrentMap.cookieMap[cookie.Value]
	//concurrentMap.RUnlock()
        //getName(cookie.Value)
        name := getName(cookie.Value, w, r)
        //name := ""
	if len(name) == 0 {
		log.Info("There is no name stored for the UUID")
		return false, ""
	} else if len(name) > 0 {
		log.Infof("User is logged in with these values: name: %s. UUID: %s", name, cookie.Value)
		return true, name
	} else {
		log.Critical("There is an unknown error")
		return false, ""
	}
}

// Will parse the form to add the user into the logged in
// users and redirect to the homepage
func loginPage(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	formName := r.FormValue("name")
	if len(formName) > 0 {
		cookie := setCookie(w, r)
		//cookie := CookieManagement.SetCookie(w, r)
                //
                fmt.Println(cookie)
                //
		//concurrentMap.Lock()
		//concurrentMap.cookieMap[cookie.Value] = formName
		//concurrentMap.Unlock()
                setName(cookie.Value, formName)
		log.Debugf("Name passed in: %s", formName)
		http.Redirect(w, r, "/", 302)
		return
	} else {
		log.Info("Error! User did not pass in a name to /login")
		renderNoNamePage(w)
		return
	}
}

// Here is the logout page that will remove the cookies assosiated with the user
func logoutPage(w http.ResponseWriter, r *http.Request) {
	printRequests(r)
	cookie := setCookie(w, r)
	//cookie := CookieManagement.SetCookie(w, r)
	//concurrentMap.Lock()
	//name := concurrentMap.cookieMap[cookie.Value]
	//delete(concurrentMap.cookieMap, cookie.Value)
	//concurrentMap.Unlock()
        setName(cookie.Value, "")
        //DELETE NAME SOMEHOW
        name := ""
	log.Debugf("Deleting %s and %s from the server", cookie.Value, name)
        deleteCookie(w, r)
	//deletingCookie := &http.Cookie{Name: "uuid", Value: "s", Expires: time.Unix(1, 0), HttpOnly: true}
	//http.SetCookie(w, deletingCookie)
	logoutTemplatesSlice := make([]string, len(templatesSlice))
	copy(logoutTemplatesSlice, templatesSlice)
	logoutTemplatesSlice = append(logoutTemplatesSlice, fmt.Sprintf("%s/logout.html", templatesFolder))
	var logoutPage = template.Must(template.New("logout").ParseFiles(logoutTemplatesSlice...))
	logoutPage.ExecuteTemplate(w, "template", "")
}

// Function for printing the request URL path
func printRequests(r *http.Request) {
	urlPath := r.URL.Path
	log.Infof("Request url path: %s", urlPath)
}

// Sets up the templates slice
func templateSetup() {
	templatesSlice = append(templatesSlice, fmt.Sprintf("%s/template.html", templatesFolder))
	templatesSlice = append(templatesSlice, fmt.Sprintf("%s/menu.html", templatesFolder))
}


// Main handler that runs the server on the port or shows the version of the server
func main() {
	defer log.Flush()
	templatesFolder = templatesFlag
	templateSetup()
	logFileName := fmt.Sprintf("etc/%s.xml", logFile)
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
	log.Infof("Port flag is set as: %d", port)
	log.Infof("Version flag is set? %v", version)
	log.Infof("Log config file flag is set as: %s", logFile)
	log.Infof("Templates folder flag is set as: %s", templatesFlag)
	log.Info("Server has started up!")
	if version {
		log.Info("Printing out the version")
		fmt.Println("Assignment Version: 3")
		return
	}
	http.HandleFunc("/time", timeHandler)
	http.HandleFunc("/index.html", indexPage)
	http.HandleFunc("/login", loginPage)
	http.HandleFunc("/logout", logoutPage)
	http.HandleFunc("/", errorer)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	var portString = fmt.Sprintf(":%d", port)
	err := http.ListenAndServe(portString, nil)
	if err != nil {
		log.Errorf("Server Failed: %s", err)
		os.Exit(1)
	}
	log.Info("Server Closed")
}
