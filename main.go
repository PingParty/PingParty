package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

var templates = template.Must(template.New("").ParseGlob("templates/*.html"))

func main() {
	r := mux.NewRouter()
	r.Path("/login").HandlerFunc(login)
	r.Path("/").HandlerFunc(homepage)

	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", r)
}

func homepage(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "lohp.html", nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "login.html", nil)
}
