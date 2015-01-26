package CookieManagement

import (
	log "../seelog-master/"
	"net/http"
	"os"
	"os/exec"
	"time"
	"sync"
)

// Stores the cookie information
var concurrentMap struct {
	sync.RWMutex
	cookieMap map[string]Person
}

// Intitalizes the concurrentMap
func init() {
	concurrentMap = struct {
		sync.RWMutex
		cookieMap map[string]Person
	}{cookieMap: make(map[string]Person)}
	log.Debug("Initalizing the map")
        //concurrentMap.cookieMap["asf"] = "tasting"
}

//type CookieManager struct {
//	Name string
//	Num  int
//}

//type Manager interface {
//  TestSetCookie(w http.ResponseWriter, r *http.Request)
//}
//func TestSetCookie(w http.ResponseWriter, r *http.Request){
//  log.Info("COOKIE MANAGER TESST")
//}

type Person struct {
  Name string
}

func GetName(s string) Person {
  return concurrentMap.cookieMap[s]
}

func SetName(uuid string, name string) {
		concurrentMap.Lock()
  concurrentMap.cookieMap[uuid] = Person{Name: name}
}

func DeletePerson(uuid string) Person {
		//concurrentMap.Lock()
		//person := concurrentMap.cookieMap[cookie.Value]
		//delete(concurrentMap.cookieMap, cookie.Value)
		//concurrentMap.Unlock()
                //return person
                return Person{Name: ""}
}

//func NewCookieManager() *CookieManager {
//	log.Info("testing from cookie Manager")
//	return &CookieManager{}
//}

//func setCookie(w http.ResponseWriter, r *http.Request) *http.Cookie {
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
	log.Infof("Setting cookie with UUID: %s", uuid)
	uuidLen := len(uuid) - 1
	uuidString := string(uuid[:uuidLen])
	cookie := &http.Cookie{Name: "uuid", Value: uuidString, Expires: time.Now().Add(356 * 24 * time.Hour), HttpOnly: true}
	http.SetCookie(w, cookie)
	return cookie
}
