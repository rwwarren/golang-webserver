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
	"../cookieManagement/"
	log "../seelog-master/"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"sync"
	"time"
)

// Stores the cookie information
var concurrentMap struct {
	sync.RWMutex
	cookieMap map[string]string
}

type PageInformation struct {
	Name        string
	CurrentTime string
	UTCtime     string
}

var templatesFolder string
var templatesSlice []string

// Intitalizes the concurrentMap
func init() {
	log.Debug("Logger not initialized yet")
	concurrentMap = struct {
		sync.RWMutex
		cookieMap map[string]string
	}{cookieMap: make(map[string]string)}
	log.Debug("Initalizing the map")
}

// Handles the timeserver which shows the current time
// for the local timezone
func timeHandler(w http.ResponseWriter, r *http.Request) {
	printRequests(r)
	const layout = "3:04:05 PM"
	const UTClayout = "15:04:05 MST"
	personalString := ""
	cookie := CookieManagement.SetCookie(w, r)
	concurrentMap.RLock()
	if len(concurrentMap.cookieMap[cookie.Value]) > 0 {
		name := concurrentMap.cookieMap[cookie.Value]
		personalString = fmt.Sprintf(", %s", name)
		log.Debugf("User is logged in as: %s", name)
	}
	concurrentMap.RUnlock()
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
	cookie := CookieManagement.SetCookie(w, r)
	concurrentMap.RLock()
	name := concurrentMap.cookieMap[cookie.Value]
	concurrentMap.RUnlock()
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
		cookie := CookieManagement.SetCookie(w, r)
		concurrentMap.Lock()
		concurrentMap.cookieMap[cookie.Value] = formName
		concurrentMap.Unlock()
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
	cookie := CookieManagement.SetCookie(w, r)
	concurrentMap.Lock()
	name := concurrentMap.cookieMap[cookie.Value]
	delete(concurrentMap.cookieMap, cookie.Value)
	concurrentMap.Unlock()
	log.Debugf("Deleting %s and %s from the server", cookie.Value, name)
	deletingCookie := &http.Cookie{Name: "uuid", Value: "s", Expires: time.Unix(1, 0), HttpOnly: true}
	http.SetCookie(w, deletingCookie)
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
	port := flag.Int("port", 8080, "Set the server port, default port: 8080")
	version := flag.Bool("V", false, "Shows the version of the timeserver")
	logFile := flag.String("log", "logConfig", "This is the logger configuration file")
	templatesFlag := flag.String("templates", "templates", "This is the templates folder name")
	flag.Parse()
	templatesFolder = *templatesFlag
	templateSetup()
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
	log.Infof("Port flag is set as: %d", *port)
	log.Infof("Version flag is set? %v", *version)
	log.Infof("Log config file flag is set as: %s", *logFile)
	log.Infof("Templates folder flag is set as: %s", *templatesFlag)
	log.Info("Server has started up!")
	if *version {
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
	var portString = fmt.Sprintf(":%d", *port)
	err := http.ListenAndServe(portString, nil)
	if err != nil {
		log.Errorf("Server Failed: %s", err)
		os.Exit(1)
	}
	log.Info("Server Closed")
}