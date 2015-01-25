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
// "-LogOutput FILENAME" pass in the log file and it will output to FILENAME.log

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"
        "text/template"
        log "../seelog-master/"
)

// Stores the cookie information
var concurrentMap struct {
	sync.RWMutex
	cookieMap map[string]string
}

// Intitalizes the concurrentMap
func init() {
	concurrentMap = struct {
		sync.RWMutex
		cookieMap map[string]string
	}{cookieMap: make(map[string]string)}
}

// Handles the timeserver which shows the current time
// for the local timezone
func timeHandler(w http.ResponseWriter, r *http.Request) {
	printRequests(r)
	const layout = "3:04:05 PM"
	const UTClayout = "15:04:05 MST"
	personalString := ""
	cookie := setCookie(w, r)
	concurrentMap.RLock()
	if len(concurrentMap.cookieMap[cookie.Value]) > 0 {
		name := concurrentMap.cookieMap[cookie.Value]
		personalString = fmt.Sprintf(", %s", name)
	}
	concurrentMap.RUnlock()
//	fmt.Fprintf(w, `<html><head><style>
//          p {font-size: xx-large}
//          span.time {color: red}
//          </style>
//          </head>
//          <body>
//          <p>The time is now <span class="time">%s</span> (%s)%s.</p>
//          </body>
//          </html>`, time.Now().Local().Format(layout),
//              time.Now().UTC().Format(UTClayout), personalString)
//	return
  //var hogeTmpl = template.New("template").ParseFiles("templates/template.html", "templates/menu.html", "templates/time.html")
  var hogeTmpl = template.Must(template.New("template").ParseFiles("templates/template.html", "templates/menu.html", "templates/time.html"))
  fmt.Println(personalString)
  currentTime := time.Now().Local().Format(layout)
  UTCTime := time.Now().UTC().Format(UTClayout)
  asdf := &Testing{
    Name: personalString,
    CurrentTime: currentTime,
    UTCtime: UTCTime,
  }
  hogeTmpl.ExecuteTemplate(w, "template", asdf)
  data := asdf
  fmt.Println(asdf)
  fmt.Println(data)
  //hogeTmpl.ExecuteTemplate(w, "template", map[string]string{"Name":personalString})
  //hogeTmpl.Execute(w, "template", data)
  //hogeTmpl.ExecuteTemplate(w, "template", data)
  //hogeTmpl.ExecuteTemplate(w, "template", *asdf)
  //hogeTmpl.ExecuteTemplate(w, "template", "personalString")
  //hogeTmpl.ExecuteTemplate(w, "template", personalString)
  //hogeTmpl.ExecuteTemplate(w, "template", []string{time.Now().Local().Format(layout), time.Now().UTC().Format(UTClayout), personalString})
  //hogeTmpl.ExecuteTemplate(w, "template", "Hoge")
}

type Testing struct {
  Name          string
  CurrentTime   string
  UTCtime       string
}

// Handles errors for when the page is not found
func errorer(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		indexPage(w, r)
		return
	}
	printRequests(r)
	log.Info("Error, url not found: These are not the URLs you are looking for.\n")
	//log.Println("Error, url not found: These are not the URLs you are looking for.")
	w.WriteHeader(404)
	fmt.Fprintf(w, `<html><head><style>
          p {font-size: xx-large}
          span.time {color: red}
          </style>
          </head>
          <body>
          <p>These are not the URLs you're looking for.</p>
          </body>
          </html>`)
	return
}

// Checks the index page and will render the
// index based on the user being loggedin or not
func indexPage(w http.ResponseWriter, r *http.Request) {
	printRequests(r)
	isLoggedIn, name := checkLogin(w, r)
	if isLoggedIn {
		renderIndex(w, name)
		return
	} else {
		renderLogin(w, r)
		return
	}

}

// Renders the page for a loggedin user
func renderIndex(w http.ResponseWriter, name string) {
//	fmt.Fprintf(w, `<html>
//      <body>
//      Greetings, %s.
//      </body>
//      </html>`, name)
//	return
//  var hogeTmpl = template.Must(template.New("hoge").ParseFiles("templates/template.html"))
//  var hogeTmpl = template.Must(template.New("hoge").ParseFiles("templates/template.html", "templates/menu.html"))
  var hogeTmpl = template.Must(template.New("hoge").ParseFiles("templates/template.html", "templates/menu.html", "templates/index.html"))
  person := &Testing{
    Name: name,
  }
  //hogeTmpl.ExecuteTemplate(w, "template", "Hogeasdfasfdasdf")
  hogeTmpl.ExecuteTemplate(w, "template", person)


}

// Renders the page if there is no name passed into
// the login page
func renderNoNamePage(w http.ResponseWriter) {
//	fmt.Fprintf(w, `<html>
//          <body>
//          C'mon, I need a name.
//          </body>
//          </html>`)
//	return
  var logoutPage = template.Must(template.New("hoge").ParseFiles("templates/template.html", "templates/menu.html", "templates/noNamePage.html"))
  logoutPage.ExecuteTemplate(w, "template", "")
}

// Renders the login page to the website
func renderLogin(w http.ResponseWriter, r *http.Request) {
//	fmt.Fprintf(w, `<html>
//          <body>
//          <form action="login">
//            What is your name, Earthling?
//            <input type="text" name="name" size="50">
//            <input type="submit">
//          </form>
//          </p>
//          </body>
//          </html>`)
//	return
  var loginPage = template.Must(template.New("hoge").ParseFiles("templates/template.html", "templates/menu.html", "templates/loginPage.html"))
  loginPage.ExecuteTemplate(w, "template", "")
}

// Returns the cookie for the server. Will set one if there is none
func setCookie(w http.ResponseWriter, r *http.Request) *http.Cookie {
	checkCookie, cookieError := r.Cookie("uuid")
	if cookieError == nil {
		log.Infof("Cookie is already set: %s", checkCookie.Value)
		//log.Printf("Cookie is already set: %s", checkCookie.Value)
		return checkCookie
	}
	uuid, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Infof("Error something went wrong with uuidgen: %s \n", err)
		//log.Printf("Error something went wrong with uuidgen: %s \n", err)
		os.Exit(1)
	}
	log.Infof("Setting cookie with UUID: %s", uuid)
	//log.Printf("Setting cookie with UUID: %s", uuid)
	uuidLen := len(uuid) - 1
	uuidString := string(uuid[:uuidLen])
	cookie := &http.Cookie{Name: "uuid", Value: uuidString, Expires: time.Now().Add(356 * 24 * time.Hour), HttpOnly: true}
	http.SetCookie(w, cookie)
	return cookie
}

// Checks the if the user is logged in and if there is a user
// associated with the cookie
func checkLogin(w http.ResponseWriter, r *http.Request) (bool, string) {
	cookie := setCookie(w, r)
	concurrentMap.RLock()
	name := concurrentMap.cookieMap[cookie.Value]
	concurrentMap.RUnlock()
	if len(name) == 0 {
		log.Info("There is no name stored for the UUID")
		//log.Println("There is no name stored for the UUID")
		return false, ""
	} else if len(name) > 0 {
		log.Infof("User is logged in with these values: name: %s. UUID: %s", name, cookie.Value)
		//log.Printf("User is logged in with these values: name: %s. UUID: %s", name, cookie.Value)
		return true, name
	} else {
		log.Info("There is an unknown error")
		//log.Println("There is an unknown error")
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
		concurrentMap.Lock()
		concurrentMap.cookieMap[cookie.Value] = formName
		concurrentMap.Unlock()
		http.Redirect(w, r, "/", 302)
		return
	} else {
		log.Info("Error! User did not pass in a name to /login")
		//log.Println("Error! User did not pass in a name to /login")
		renderNoNamePage(w)
		return
	}
}

// Here is the logout page that will remove the cookies assosiated with the user
func logoutPage(w http.ResponseWriter, r *http.Request) {
	printRequests(r)
	if cookie, err := r.Cookie("uuid"); err == nil {
		concurrentMap.Lock()
		name := concurrentMap.cookieMap[cookie.Value]
		delete(concurrentMap.cookieMap, cookie.Value)
		concurrentMap.Unlock()
		log.Infof("Deleting %s and %s from the server\n", cookie.Value, name)
		//log.Printf("Deleting %s and %s from the server\n", cookie.Value, name)
	}
	cookie := &http.Cookie{Name: "uuid", Value: "s", Expires: time.Unix(1, 0), HttpOnly: true}
	http.SetCookie(w, cookie)
//	fmt.Fprintf(w, `<html>
//          <head>
//          <META http-equiv="refresh" content="10;URL=/">
//          <body>
//          <p>Good-bye.</p>
//          </body>
//          </html>`)
//	return
  var logoutPage = template.Must(template.New("hoge").ParseFiles("templates/template.html", "templates/menu.html", "templates/logout.html"))
  logoutPage.ExecuteTemplate(w, "template", "")
}

// Function for printing the request URL path
func printRequests(r *http.Request) {
	urlPath := r.URL.Path
	log.Infof("Request url path: %s \n", urlPath)
	//log.Printf("Request url path: %s \n", urlPath)
}

// Main handler that runs the server on the port or shows the version of the server
func main() {
        defer log.Flush()
	port := flag.Int("port", 8080, "Set the server port, default port: 8080")
	version := flag.Bool("V", false, "Shows the version of the timeserver")
	logFile := flag.String("LogOutput", "", "This is the log output file name")
	flag.Parse()
        logger, logError := log.LoggerFromConfigAsFile("logConfig.xml")
        if logError != nil {
              fmt.Printf("Log instantiation error: %s", logError)
        }
        log.ReplaceLogger(logger)
	//if len(*logFile) > 0 {
	//	logFileName := fmt.Sprintf("%s.log", *logFile)
	//	f, logerr := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	//	if logerr != nil {
	//		fmt.Printf("Error opening the log file: %v", logerr)
	//		os.Exit(1)
	//	}
	//	defer f.Close()
	//	log.SetOutput(f)
	//}
	log.Infof("Port flag is set as: %d\n", *port)
	log.Infof("Version flag is set? %v\n", *version)
	log.Infof("Log file flag is set as: %s\n", *logFile)
	log.Info("Server has started up!")
	//log.Printf("Port flag is set as: %d\n", *port)
	//log.Printf("Version flag is set? %v\n", *version)
	//log.Printf("Log file flag is set as: %s\n", *logFile)
	//log.Println("Server has started up!")
	if *version {
		log.Info("Printing out the version")
		//log.Println("Printing out the version")
		fmt.Println("Assignment Version: 2")
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
		log.Infof("Server Failed: %s\n", err)
		//log.Printf("Server Failed: %s\n", err)
		os.Exit(1)
	}
	log.Info("Server Closed")
	//log.Println("Server Closed")
}
