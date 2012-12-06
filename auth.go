package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/surma-dump/goauth2/oauth"
)

type Authenticator interface {
	http.Handler
}

type Extractor interface {
	Extract(c *http.Client) (string, error)
}

type ExtractorFunc func(c *http.Client) (string, error)

func (e ExtractorFunc) Extract(c *http.Client) (string, error) {
	return e(c)
}

type OAuthAuthenticator struct {
	config    *oauth.Config
	extractor Extractor
	*mux.Router
}

func NewOAuthAuthenticator(c *oauth.Config, e Extractor) *OAuthAuthenticator {
	a := &OAuthAuthenticator{
		config:    c,
		extractor: e,
	}
	a.Router = mux.NewRouter()
	a.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		a.authHandler(w, r)
	})
	a.Router.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		a.authCallbackHandler(w, r)
	})
	return a
}

func (a *OAuthAuthenticator) authHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, a.config.AuthCodeURL(""), http.StatusFound)
}

func (a *OAuthAuthenticator) authCallbackHandler(w http.ResponseWriter, r *http.Request) {
	transport := (&oauth.Transport{Config: a.config})
	_, err := transport.Exchange(r.FormValue("code"))
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not get access token: %s", err), http.StatusServiceUnavailable)
		return
	}

	uid, err := a.extractor.Extract(transport.Client())
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not get user id: %s", err), http.StatusServiceUnavailable)
		return
	}
	fmt.Fprintf(w, "ID: %s", uid)
}

func GitHubExtractor(c *http.Client) (string, error) {
	r, err := c.Get("https://api.github.com/user")
	if err != nil {
		return "", err
	}

	user := struct {
		Id int `json:"id"`
	}{
		Id: -1,
	}
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return "", err
	}
	if user.Id == -1 {
		err = fmt.Errorf("Invalid user id")
	}
	return fmt.Sprintf("%d", user.Id), err
}

func GoogleExtractor(c *http.Client) (string, error) {
	r, err := c.Get("https://www.googleapis.com/oauth2/v1/userinfo")
	if err != nil {
		return "", err
	}

	user := struct {
		Email string `json:"email"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return "", err
	}
	if user.Email == "" {
		err = fmt.Errorf("Invalid user id")
	}
	return user.Email, err
}
