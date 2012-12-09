package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/surma-dump/mux"

	"code.google.com/p/goauth2/oauth"
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
	a.Router.KeepContext = true
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
}

func NewJSONExtractor(url string, field string) Extractor {
	return ExtractorFunc(func(c *http.Client) (string, error) {
		r, err := c.Get(url)
		if err != nil {
			return "", err
		}

		var data interface{}
		err = json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			return "", err
		}
		switch x := data.(type) {
		case map[string]interface{}:
			data = x[field]
		default:
			return "", fmt.Errorf("Unhandled JSON type")
		}
		switch x := data.(type) {
		case float64:
			return fmt.Sprintf("%.0f", data), nil
		case string:
			return x, nil
		}
		return "", fmt.Errorf("Unsupported id type")
	})
}

type contextKey int

var (
	sessionKey contextKey = 0
)

func NewSessionOpener(s sessions.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.Get(r, "uid")
		if err != nil {
			session, err = s.New(r, "uid")
			if err != nil {
				http.Error(w, fmt.Sprintf("Could not create session: %s", err), http.StatusInternalServerError)
			}
		}
		context.Set(r, sessionKey, session)
	})
}

func SessionSaver(w http.ResponseWriter, r *http.Request) {
	err := context.Get(r, sessionKey).(*sessions.Session).Save(r, w)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not save session: %s", err), http.StatusInternalServerError)
	}
}
