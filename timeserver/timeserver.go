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
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
	"../cookieManagement/"
	log "../seelog-master/"
	"html/template"
)

// Stores the cookie information
var concurrentMap struct {
	sync.RWMutex
	cookieMap map[string]string
}

//TODO rename this
type Testing struct {
	Name        string
	CurrentTime string
	UTCtime     string
}

var templatesFolder string

// Intitalizes the concurrentMap
func init() {
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
	var timeTmpl = template.Must(template.New("time").ParseFiles("templates/template.html", "templates/menu.html", "templates/time.html"))
	currentTime := time.Now().Local().Format(layout)
	UTCTime := time.Now().UTC().Format(UTClayout)
	data := &Testing{
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
	var errorPage = template.Must(template.New("ErrorPage").ParseFiles("templates/template.html", "templates/menu.html", "templates/404.html"))
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
	var indexPage = template.Must(template.New("hoge").ParseFiles("templates/template.html", "templates/menu.html", "templates/index.html"))
	person := &Testing{
		Name: name,
	}
	indexPage.ExecuteTemplate(w, "template", person)
}

// Renders the page if there is no name passed into
// the login page
func renderNoNamePage(w http.ResponseWriter) {
	var logoutPage = template.Must(template.New("NoNamePage").ParseFiles("templates/template.html", "templates/menu.html", "templates/noNamePage.html"))
	logoutPage.ExecuteTemplate(w, "template", "")
}

// Renders the login page to the website
func renderLogin(w http.ResponseWriter, r *http.Request) {
	var loginPage = template.Must(template.New("Login").ParseFiles("templates/template.html", "templates/menu.html", "templates/loginPage.html"))
	loginPage.ExecuteTemplate(w, "template", "")
}

//// Returns the cookie for the server. Will set one if there is none
//func setCookie(w http.ResponseWriter, r *http.Request) *http.Cookie {
//	checkCookie, cookieError := r.Cookie("uuid")
//	if cookieError == nil {
//		log.Infof("Cookie is already set: %s", checkCookie.Value)
//		return checkCookie
//	}
//	uuid, err := exec.Command("uuidgen").Output()
//	if err != nil {
//		log.Criticalf("Error something went wrong with uuidgen: %s", err)
//		os.Exit(1)
//	}
//	log.Infof("Setting cookie with UUID: %s", uuid)
//	uuidLen := len(uuid) - 1
//	uuidString := string(uuid[:uuidLen])
//	cookie := &http.Cookie{Name: "uuid", Value: uuidString, Expires: time.Now().Add(356 * 24 * time.Hour), HttpOnly: true}
//	http.SetCookie(w, cookie)
//	return cookie
//}

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
        cookie := CookieManagement.SetCookie(w,r)
	concurrentMap.Lock()
	name := concurrentMap.cookieMap[cookie.Value]
	delete(concurrentMap.cookieMap, cookie.Value)
	concurrentMap.Unlock()
	log.Debugf("Deleting %s and %s from the server", cookie.Value, name)
	deletingCookie := &http.Cookie{Name: "uuid", Value: "s", Expires: time.Unix(1, 0), HttpOnly: true}
	http.SetCookie(w, deletingCookie)
	var logoutPage = template.Must(template.New("logout").ParseFiles("templates/template.html", "templates/menu.html", "templates/logout.html"))
	logoutPage.ExecuteTemplate(w, "template", "")
}

// Function for printing the request URL path
func printRequests(r *http.Request) {
	urlPath := r.URL.Path
	log.Infof("Request url path: %s", urlPath)
}

// Main handler that runs the server on the port or shows the version of the server
func main() {
	defer log.Flush()
	port := flag.Int("port", 8080, "Set the server port, default port: 8080")
	version := flag.Bool("V", false, "Shows the version of the timeserver")
	logFile := flag.String("log", "logConfig", "This is the logger configuration file")
	templatesFlag := flag.String("templates", "templates/", "This is the templates folder name")
	flag.Parse()
	templatesFolder = *templatesFlag
	logFileName := fmt.Sprintf("etc/%s.xml", *logFile)
	logger, logError := log.LoggerFromConfigAsFile(logFileName)
	if logError != nil {
		fmt.Printf("Log instantiation error: %s", logError)
	}
	log.ReplaceLogger(logger)
	log.Infof("Port flag is set as: %d", *port)
	log.Infof("Version flag is set? %v", *version)
	log.Infof("Log file flag is set as: %s", *logFile)
	log.Info("Server has started up!")
	if *version {
		log.Info("Printing out the version")
		fmt.Println("Assignment Version: 2")
		return
	}
	//cookieManager := CookieManagement.NewCookieManager()
        //cookieManager.SetCookie()
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
	//cookieManager := CookieManagement.NewCookieManager()
	//cookieManager := NewCookieManager()
	//cookieManager := new(CookieManager)
	//fmt.Println(cookieManager)
	log.Info("Server Closed")
}
