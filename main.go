package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/PingParty/PingParty/db"
	"github.com/PingParty/PingParty/models"
	"github.com/gorilla/mux"

	_ "github.com/joho/godotenv/autoload" // load .env into local environment if it exists
)

var devMode = os.Getenv("DEVMODE") == "true"

var templates = template.Must(template.New("").ParseGlob("templates/*.html"))
var data *db.DB

func main() {
	var err error
	if data, err = db.New(os.Getenv("DB_CONN")); err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()

	// TODO: dedicated login landing page r.Path("/login")
	r.Path("/auth/github").HandlerFunc(redirectToGithub)
	r.Path("/auth/github/callback").HandlerFunc(githubCallback)

	r.Path("/signup").Methods(http.MethodGet).HandlerFunc(signupForm)
	r.Path("/signup").Methods(http.MethodPost).HandlerFunc(completeSignup)

	r.Path("/").HandlerFunc(homepage)

	fmt.Println("Listening on port 8000")
	http.ListenAndServe(":8000", r)
}

func render(w http.ResponseWriter, name string, ctx interface{}) {
	if devMode {
		templates = template.Must(template.New("").ParseGlob("templates/*.html"))
	}
	templates.ExecuteTemplate(w, name, ctx)
}

func homepage(w http.ResponseWriter, r *http.Request) {
	render(w, "lohp.html", nil)
}

func signupForm(w http.ResponseWriter, r *http.Request) {
	sess, _ := cookieStore.Get(r, cookieNameAuth)
	// TODO: handle session error
	u, ok := sess.Values["user"].(*models.User)
	if !ok {
		errPage(w, "No login found. Please go to /login")
		return
	}
	render(w, "signup.html", u)
}

func completeSignup(w http.ResponseWriter, r *http.Request) {
	sess, _ := cookieStore.Get(r, cookieNameAuth)
	// TODO: handle session error
	u, ok := sess.Values["user"].(*models.User)
	if !ok {
		errPage(w, "No login found. Please go to /login")
		return
	}
	r.ParseForm()
	fmt.Println(u, r.FormValue("email"))
}

func errPage(w http.ResponseWriter, msg string) {
	http.Error(w, msg, 500)
}
