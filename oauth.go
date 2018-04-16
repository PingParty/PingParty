package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/PingParty/PingParty/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var githubConfig = &oauth2.Config{
	ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
	ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
	Endpoint:     github.Endpoint,
	Scopes:       []string{"user:email"},
}

func redirectToGithub(w http.ResponseWriter, r *http.Request) {
	state := "abc" // TODO: random state. Stored in cookie
	http.Redirect(w, r, githubConfig.AuthCodeURL(state), http.StatusTemporaryRedirect)
}

func githubCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	// Get the state and code from the url parameters
	if state != "abc" { //todo: compare with state stored in cookie
		http.Error(w, "Oauth Error", 400)
		return
	}
	if code == "" {
		http.Error(w, "Oauth Error", 400)
		return
	}
	// Now we go back to github and exchange the code for a valid oauth token.
	token, err := githubConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Oauth Error", 400)
		return
	}
	// Cool. We have a valid github token. Now let's learn about the github user they are:
	ghInfo, err := getGithubInfo(token)
	if err != nil {
		http.Error(w, "Oauth Error", 400)
		return
	}
	// Awesome. Now we know who they are on github. Let's see if we have already created a user for them:
	user, err := data.GetExistingThirdPartyUser("github", fmt.Sprint(ghInfo.ID))
	if err != nil {
		http.Error(w, "Oauth Error", 400)
		return
	}
	// We know them! Set a cookie and send them along!
	if user != nil {
		setSessionCookie(user, w)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect) // todo: get destination from cookie
		return
	}
	// Never seen them before. Let's get some information from them to get the email they want to use, and we can create an account.
	// First, let's get their emails from github:
	emails, err := getGithubEmails(token)
	if err != nil {
		http.Error(w, "Oauth Error", 400)
		fmt.Println(err)
		return
	}
	// Store partial user info in a secure cookie, and show them a "choose email" page
	user = &models.User{
		LoginType: "github",
		LoginID:   fmt.Sprint(ghInfo.ID),
		LoginName: &ghInfo.Login,
	}
	setCookie("signupUser", user, 10, w)
	templates.ExecuteTemplate(w, "createUser.html", emails)
}

type githubLoginInfo struct {
	Login string `json:"login"`
	ID    int    `json:"id"`
	Email string `json:"email"`
}
type githubEmail struct {
	Email   string `json: "email"`
	Primary bool   `json:"primary"`
}

func getGithubInfo(token *oauth2.Token) (*githubLoginInfo, error) {
	info := &githubLoginInfo{}
	if err := githubGet(token, "https://api.github.com/user", info); err != nil {
		return nil, err
	}
	return info, nil
}

func githubGet(token *oauth2.Token, url string, dest interface{}) error {
	r, _ := http.NewRequest(http.MethodGet, url, nil)
	token.SetAuthHeader(r)
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	return dec.Decode(dest)
}

func getGithubEmails(token *oauth2.Token) ([]githubEmail, error) {
	emails := []githubEmail{}
	if err := githubGet(token, "https://api.github.com/user/emails", &emails); err != nil {
		return nil, err
	}
	return emails, nil
}