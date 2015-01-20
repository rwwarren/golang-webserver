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

package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"
        "os"
        "os/exec"
        //"sync"
        "sync"
        "log"
)

var cookieMap map[string]string
var counter struct{
      sync.RWMutex
      m map[string]string
      }

func init() {
  counter = struct{
        sync.RWMutex
        m map[string]string
        }{m: make(map[string]string)}

  cookieMap = make(map[string]string)
}

// Handles the timeserver which shows the current time
// for the local timezone
func timeHandler(w http.ResponseWriter, r *http.Request) {
        printRequests(r)
        //TODO if user logged say greetings
	const layout = "3:04:05 PM"
        personalString := ""
        //TODO fix this
        cookie, err := r.Cookie("uuid")
        counter.RLock()
        if err == nil && len(cookieMap[cookie.Value]) > 0 {
          name := cookieMap[cookie.Value]
          personalString = fmt.Sprintf(", %s", name)
        }
        counter.RUnlock()
	fmt.Fprintf(w, `<html><head><style>
          p {font-size: xx-large}
          span.time {color: red}
          </style>
          </head>
          <body>
          <p>The time is now <span class="time">%s</span>%s.</p>
          </body>
          </html>`, time.Now().Local().Format(layout), personalString)
        return
}

// Handles errors for when the page is not found
func errorer(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/" {
          indexPage(w, r)
          //loginForm(w, r)
          return
        }
        printRequests(r)
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

func indexPage(w http.ResponseWriter, r *http.Request) {
  isLoggedIn, name:= checkLogin(w, r)
  if isLoggedIn {
    //something
    renderIndex(w, name)
    //renderIndex(w, r)
    return
  } else {
    //please login
    renderLogin(w, r)
    return
  }

}

func renderIndex(w http.ResponseWriter, name string){
    fmt.Fprintf(w, `<html>
      <body>
      Greetings, %s.
      </body>
      </html>`, name)
    return
}

func renderLogin(w http.ResponseWriter, r *http.Request) {
          r.ParseForm()
          formName := r.FormValue("name")
          log.Printf("Form request information: %s\n", formName)
          if len(formName) > 0 {
              loginPage(w, r)
              return
          }
          setCookie(w)
	fmt.Fprintf(w, `<html>
          <body>
          <form action="login">
            What is your name, Earthling?
            <input type="text" name="name" size="50">
            <input type="submit">
          </form>
          </p>
          </body>
          </html>`)
          return
}

func renderNoNamePage(w http.ResponseWriter) {
    fmt.Fprintf(w, `<html>
          <body>
          C'mon, I need a name.
          </body>
          </html>`)
    return
}

func setCookie(w http.ResponseWriter) {
          uuid, err := exec.Command("uuidgen").Output()
          if err != nil {
                log.Printf("Error something went wrong with uuidgen: %s \n", err)
                os.Exit(1)
          }
          log.Printf("Setting cookie with UUID: %s", uuid)
          uuidLen := len(uuid)-1
          uuidString := string(uuid[:uuidLen])
          cookie := &http.Cookie{Name:"uuid", Value:uuidString, Expires:time.Now().Add(356*24*time.Hour), HttpOnly:true}
          http.SetCookie(w, cookie)
          return
}


func checkLogin(w http.ResponseWriter, r *http.Request) (bool, string) {
  cookie, err := r.Cookie("uuid")
  if err != nil {
    log.Printf("Error with the cookie: %s", err)
    return false, ""
  } else if len(cookie.Value) == 0 {
    log.Println("There is no cookie, currently setting the cookie")
    setCookie(w)
    //Set the cookie
    return false, ""
  } else if len(cookieMap[cookie.Value]) == 0 {
    log.Println("There is no name stored for the UUID")
    return false, ""
  } else if len(cookieMap[cookie.Value]) > 0 {
    //We have a valid uuid & cookie
    log.Printf("User is logged in with these values: name: %s. UUID: %s", cookieMap[cookie.Value], cookie.Value)
    return true, cookieMap[cookie.Value]
  } else {
    log.Println("There is an unknown error")
    return false, ""
  }

}

func loginPage(w http.ResponseWriter, r *http.Request){
        r.ParseForm()
        formName := r.FormValue("name")
        if len(formName) > 0 {
            cookie, err := r.Cookie("uuid")
            if err != nil {
              log.Printf("error getting the cookie: %s", err)
              os.Exit(1)
            }
            cookieMap[cookie.Value] = formName
            http.Redirect(w, r, "/", 302)
            return
        } else {
            log.Println("Error! User did not pass in a name to /login")
            renderNoNamePage(w)
            return
        }
}

// Here is the logout page that will remove the cookies assosiated with the user
func logoutPage(w http.ResponseWriter, r *http.Request) {
        printRequests(r)
        if cookie, err := r.Cookie("uuid"); err == nil {
            log.Printf("Deleting %s and %s from the server\n", cookie.Value, cookieMap[cookie.Value])
            delete(cookieMap, cookie.Value)
        }
        cookie := &http.Cookie{Name:"uuid", Value:"s", Expires:time.Unix(1, 0), HttpOnly:true}
        http.SetCookie(w, cookie)
	fmt.Fprintf(w, `<html>
          <head>
          <META http-equiv="refresh" content="10;URL=/">
          <body>
          <p>Good-bye.</p>
          </body>
          </html>`)
        return
}

// Funtion for printing the request URL path
func printRequests(r *http.Request){
  urlPath := r.URL.Path
  log.Printf("Here is the request url path: %s \n", urlPath)
}


// Main handler that runs the server on the port or shows the version of the server
func main() {
	port := flag.Int("port", 8080, "Set the server port, default port: 8080")
	version := flag.Bool("V", false, "Shows the version of the timeserver")
	logFile := flag.String("LogOutput", "", "This is the log output file name")
	flag.Parse()
        if len(*logFile) > 0 {
          logFileName := fmt.Sprintf("%s.log", *logFile)
          f, logerr := os.OpenFile(logFileName, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
          if logerr != nil {
                fmt.Printf("Error opening the log file: %v", logerr)
                os.Exit(1)
          }
          defer f.Close()
          log.SetOutput(f)
        }
        log.Printf("Port flag is set as: %d\n", *port)
        log.Printf("Version flag is set? %v\n", *version)
        log.Printf("Log file flag is set as: %s\n", *logFile)
        log.Println("Server has started up!")
	if *version {
                log.Println("Printing out the version")
	        fmt.Println("Assignment Version: 2")
		return
	}
	http.HandleFunc("/time", timeHandler)
	http.HandleFunc("/index.html", indexPage)
	http.HandleFunc("/login", loginPage)
	http.HandleFunc("/logout", logoutPage)
	http.HandleFunc("/", errorer)
	var portString = fmt.Sprintf(":%d", *port)
	err := http.ListenAndServe(portString, nil)
        if err != nil {
	      log.Printf("Server Failed: %s\n", err)
              os.Exit(1)
        }
        log.Println("Server Closed")
}

