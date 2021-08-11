package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var (
	githubOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/callback",
		ClientID:     "Iv1.d6e26d123d63e218",
		ClientSecret: "f6bcb8ff88a127aaa7eca10a8eaacbfdd0dd6361",
		Scopes:       []string{"https://github.com/apps/goexampleapp"},
		Endpoint:     github.Endpoint,
	}
)

func main() {
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/callback", handleCallback)
	http.ListenAndServe(":8080", nil)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<html><body><a href="\login">Github Login </a></body></html>`)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	url := githubOauthConfig.AuthCodeURL("random")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != "random" {
		fmt.Println("state is not valid")
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		return
	}

	token, err := githubOauthConfig.Exchange(oauth2.NoContext, r.FormValue("code"))
	err = r.ParseForm()
	if err != nil {
		fmt.Fprintf(os.Stdout, "could not parse query: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	code := r.FormValue("code")

	reqURL := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s", "Iv1.d6e26d123d63e218", "f6bcb8ff88a127aaa7eca10a8eaacbfdd0dd6361", code)
	resp, err := http.NewRequest("POST", reqURL, nil)
	resp.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, error := client.Do(resp)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)

	data, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	data.Header.Set("Content-Type", "application/json")
	data.Header.Set("Authorization", "token "+token.AccessToken)

	client = &http.Client{}
	response, error = client.Do(data)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))

	var url = fmt.Sprintf("http://localhost:8080")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	w.WriteHeader(http.StatusFound)
}
