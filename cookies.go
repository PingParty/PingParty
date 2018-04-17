package main

import (
	"crypto/sha256"
	"log"
	"net/http"
	"os"

	"github.com/PingParty/PingParty/models"
	"github.com/gorilla/securecookie"
)

var sc *securecookie.SecureCookie

func init() {
	secret := []byte(os.Getenv("SECRET"))
	hashKey := sha256.Sum256(secret)
	blockKey := sha256.Sum256(hashKey[:])
	sc = securecookie.New(hashKey[:], blockKey[:])
	sc.SetSerializer(securecookie.JSONEncoder{})
}

func setSessionCookie(u *models.User, w http.ResponseWriter) {

}

const shortTTL = 10 * 60

func setShortCookie(name string, data interface{}, w http.ResponseWriter) {
	if encoded, err := sc.MaxAge(shortTTL).Encode(name, data); err == nil {
		cookie := &http.Cookie{
			Name:     name,
			Value:    encoded,
			Path:     "/",
			Secure:   !devMode,
			MaxAge:   10 * 60,
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
	} else {
		log.Fatal(err)
	}
}

func getCookie(name string, data interface{}, r *http.Request) error {
	c, err := r.Cookie(name)
	if err != nil {
		return err
	}
	return sc.MaxAge(shortTTL).Decode(name, c.Value, data)
}
