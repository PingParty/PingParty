package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/joho/godotenv/autoload" // load .env into local environment if it exists
)

var templates = template.Must(template.New("").ParseGlob("templates/*.html"))

func main() {
	r := mux.NewRouter()

	r.Path("/auth/github").HandlerFunc(redirectToGithub)
	r.Path("/auth/github/callback").HandlerFunc(githubCallback)

	r.Path("/").HandlerFunc(homepage)

	fmt.Println("Listening on port 8000")
	http.ListenAndServe(":8000", r)
}

func homepage(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "lohp.html", nil)
}
