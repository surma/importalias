package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/surma-dump/context"

	oauth2 "code.google.com/p/goauth2/oauth"
	oauth1 "github.com/surma-dump/oauth1a"
)

func init() {
	// Register *gouuid.UUID as a type with gob
	// as gob is being used bei gorilla’s cookie session store which
	// is used by us to save the user’s uuid.
	gob.Register(&oauth1.UserConfig{})
}

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

type OAuth1AuthenticationProvider struct {
	authname  string
	config    *oauth1.Service
	extractor Extractor
	*mux.Router
	session sessions.Store
}

func NewOAuth1AuthenticationService(name string, c *oauth1.Service, e Extractor, s sessions.Store) *OAuth1AuthenticationProvider {
	a := &OAuth1AuthenticationProvider{
		authname:  name,
		config:    c,
		extractor: e,
		session:   s,
	}
	c.Signer = &oauth1.HmacSha1Signer{}
	a.Router = mux.NewRouter()
	a.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		a.authHandler(w, r)
	})
	a.Router.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		a.authCallbackHandler(w, r)
	})
	return a
}

func (a *OAuth1AuthenticationProvider) authHandler(w http.ResponseWriter, r *http.Request) {
	uc := &oauth1.UserConfig{}
	err := uc.GetRequestToken(a.config, http.DefaultClient)
	if err != nil {
		log.Printf("Could not get OAuth1 request token for %s: %s", a.authname, err)
		// TODO: Abort gracefully
		return
	}
	url, err := uc.GetAuthorizeURL(a.config)
	if err != nil {
		log.Printf("Could not generate OAuth1 authorize url for %s: %s", a.authname, err)
		// TODO: Abort gracefully
		return
	}

	session, err := a.session.Get(r, "oauth1")
	if err != nil {
		session.Values = make(map[interface{}]interface{})
	}
	session.Values["userconfig"] = uc
	if err := session.Save(r, w); err != nil {
		http.Error(w, fmt.Sprintf("Could not save session: %s", err), http.StatusInternalServerError)
	}
	http.Redirect(w, r, url, http.StatusFound)
}

type oauth1SigningTransport struct {
	http.RoundTripper
	service    *oauth1.Service
	userconfig *oauth1.UserConfig
}

func (t *oauth1SigningTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	err := t.service.Sign(r, t.userconfig)
	if err != nil {
		return nil, err
	}
	return t.RoundTripper.RoundTrip(r)
}

func (a *OAuth1AuthenticationProvider) authCallbackHandler(w http.ResponseWriter, r *http.Request) {
	session, err := a.session.Get(r, "oauth1")
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not reopen OAuth1 %s session: %s", a.authname, err), http.StatusInternalServerError)
	}
	uc := session.Values["userconfig"].(*oauth1.UserConfig)
	token, verifier, err := uc.ParseAuthorize(r, a.config)
	if err != nil {
		log.Printf("%s: Could not get request token: %s", a.authname, err)
		http.Error(w, "Could not get request token", http.StatusServiceUnavailable)
		return
	}

	client := &http.Client{
		Transport: http.DefaultTransport,
	}
	err = uc.GetAccessToken(token, verifier, a.config, client)
	if err != nil {
		log.Printf("%s: Could not get access token: %s", a.authname, err)
		http.Error(w, "Could not get access token", http.StatusServiceUnavailable)
		return
	}
	client.Transport = &oauth1SigningTransport{http.DefaultTransport, a.config, uc}
	uid, err := a.extractor.Extract(client)
	if err != nil {
		log.Printf("%s: Could not get user id: %s", a.authname, err)
		http.Error(w, "Could not get user id", http.StatusServiceUnavailable)
		return
	}
	session.Values = make(map[interface{}]interface{})
	session.Save(r, w)
	context.Set(r, "authname", a.authname)
	context.Set(r, "authuid", uid)
}
