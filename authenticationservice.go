package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/surma-dump/context"

	oauth2 "code.google.com/p/goauth2/oauth"
)

type AuthenticationService http.Handler

type OAuth2AuthenticationProvider struct {
	authname  string
	config    *oauth2.Config
	extractor Extractor
	*mux.Router
}

func NewOAuth2AuthenticationService(name string, c *oauth2.Config, e Extractor) *OAuth2AuthenticationProvider {
	a := &OAuth2AuthenticationProvider{
		authname:  name,
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

func (a *OAuth2AuthenticationProvider) authHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, a.config.AuthCodeURL(""), http.StatusFound)
}

func (a *OAuth2AuthenticationProvider) authCallbackHandler(w http.ResponseWriter, r *http.Request) {
	transport := (&oauth2.Transport{Config: a.config})
	_, err := transport.Exchange(r.FormValue("code"))
	if err != nil {
		log.Printf("%s: Could not get access token: %s", a.authname, err)
		http.Error(w, "Could not get access token", http.StatusServiceUnavailable)
		return
	}

	uid, err := a.extractor.Extract(transport.Client())
	if err != nil {
		log.Printf("%s: Could not get user id: %s", a.authname, err)
		http.Error(w, "Could not get user id", http.StatusServiceUnavailable)
		return
	}
	context.Set(r, "authname", a.authname)
	context.Set(r, "authuid", uid)
}
