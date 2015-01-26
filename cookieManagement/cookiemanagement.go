package CookieManagement

import (
	log "../seelog-master/"
	"net/http"
	"os"
	"os/exec"
	"time"
)

type CookieManager struct {
	Name string
	Num  int
}

func NewCookieManager() *CookieManager {
	log.Info("testing from cookie Manager")
	return &CookieManager{}
}

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
	log.Infof("Setting cookie with UUID: %s", uuid)
	uuidLen := len(uuid) - 1
	uuidString := string(uuid[:uuidLen])
	cookie := &http.Cookie{Name: "uuid", Value: uuidString, Expires: time.Now().Add(356 * 24 * time.Hour), HttpOnly: true}
	http.SetCookie(w, cookie)
	return cookie
}
