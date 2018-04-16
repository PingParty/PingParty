package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

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

	if state != "abc" { //todo: compare with state stored in cookie
		// todo: fail better!
		http.Error(w, "Oauth Error", 400)
		return
	}
	if code == "" {
		// todo: fail better!
		http.Error(w, "Oauth Error", 400)
		return
	}
	token, err := githubConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Oauth Error", 400)
		return
	}
	fmt.Println(getGithubInfo(token))
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

type githubLoginInfo struct {
	Login string `json:"login"`
	Email string `json:"email"`
}

func getGithubInfo(token *oauth2.Token) (*githubLoginInfo, error) {
	r, _ := http.NewRequest(http.MethodGet, "https://api.github.com/user", nil)
	token.SetAuthHeader(r)
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	info := &githubLoginInfo{}
	if err = dec.Decode(info); err != nil {
		return nil, err
	}
	return info, nil
}
