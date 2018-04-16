package main

import (
	"crypto/sha256"
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
func setCookie(name string, data interface{}, ttlMins int, w http.ResponseWriter) {
	if encoded, err := sc.MaxAge(ttlMins*60).Encode(name, data); err == nil {
		cookie := &http.Cookie{
			Name:     name,
			Value:    encoded,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
	}
}
