package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/PingParty/PingParty/db"
	"github.com/gorilla/mux"

	_ "github.com/joho/godotenv/autoload" // load .env into local environment if it exists
)

var templates = template.Must(template.New("").ParseGlob("templates/*.html"))
var data *db.DB

func main() {
	var err error
	if data, err = db.New(os.Getenv("DB_CONN")); err != nil {
		log.Fatal(err)
	}

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
