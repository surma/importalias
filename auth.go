package main

import (
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/surma-dump/context"
	"github.com/surma-dump/gouuid"

	"code.google.com/p/goauth2/oauth"
)

func init() {
	// Register *gouuid.UUID as a type with gob
	// as gob is being used bei gorilla’s cookie session store which
	// is used by us to save the user’s uuid.
	uuid := gouuid.New()
	gob.Register(&uuid)
}

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
	authname  string
	config    *oauth.Config
	extractor Extractor
	*mux.Router
	usermgr UserManager
}

func NewOAuthAuthenticator(name string, c *oauth.Config, e Extractor, um UserManager) *OAuthAuthenticator {
	a := &OAuthAuthenticator{
		authname:  name,
		usermgr:   um,
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

var (
	CALLBACK_TEMPLATE = template.Must(template.New("").Parse(`
		<script>
			window.opener.postMessage("auth_done", "*");
			window.close();
		</script>`))
)

func (a *OAuthAuthenticator) authCallbackHandler(w http.ResponseWriter, r *http.Request) {
	transport := (&oauth.Transport{Config: a.config})
	_, err := transport.Exchange(r.FormValue("code"))
	if err != nil {
		log.Printf("%s: Could not get access token: %s", a.authname, err)
		http.Error(w, "Could not get access token", http.StatusServiceUnavailable)
		return
	}

	id, err := a.extractor.Extract(transport.Client())
	if err != nil {
		log.Printf("%s: Could not get user id: %s", a.authname, err)
		http.Error(w, "Could not get user id", http.StatusServiceUnavailable)
		return
	}

	uid, ok := context.Get(r, "uid").(*gouuid.UUID)
	// Already authenticated, add new authenticator
	if ok && uid != nil {
		err := a.usermgr.AddAuthenticator(uid, a.authname, id)
		if err != nil {
			log.Printf("Creating user failed: %s", err)
			http.Error(w, "Could not create user", http.StatusInternalServerError)
			return
		}
	} else {
		// Login
		user, err := a.usermgr.FindByAuthenticator(a.authname, id)
		if err != nil && err != ErrNotFound {
			log.Printf("Could not query user database: %s", err)
			http.Error(w, "Could not query user database", http.StatusInternalServerError)
			return
		} else if err == ErrNotFound {
			// New user
			user, err = a.usermgr.New(a.authname, id)
			if err != nil {
				log.Printf("Error creating user: %s", err)
				http.Error(w, "Error creating user", http.StatusInternalServerError)
				return
			}
		}
		session := context.Get(r, "session").(*sessions.Session)
		session.Values["uid"] = user.UID
		session.Save(r, w)
	}
	CALLBACK_TEMPLATE.Execute(w, options.Hostname)
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

func SessionHandler(s sessions.Store, ttl int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.Get(r, "session")
		if err != nil {
			session.Values = make(map[interface{}]interface{})
		}
		if uid, ok := session.Values["uid"]; ok {
			context.Set(r, "uid", uid)
		}
		session.Options.MaxAge = ttl
		context.Set(r, "session", session)
		if err := session.Save(r, w); err != nil {
			http.Error(w, fmt.Sprintf("Could not save session: %s", err), http.StatusInternalServerError)
		}
	})
}

func ValidateUID(umgr UserManager) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uid, ok := context.Get(r, "uid").(*gouuid.UUID)
		if !ok {
			http.NotFound(w, r)
			return
		}
		_, err := umgr.FindByUID(uid)
		if err != nil {
			http.NotFound(w, r)
			return
		}
	})
}

func BasicAuth(umgr UserManager) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authhdr := r.Header.Get("Authorization")
		if len(authhdr) == 0 {
			return
		}
		authhdrs := strings.Fields(authhdr)
		if len(authhdrs) != 2 || authhdrs[0] != "Basic" {
			http.NotFound(w, r)
			return
		}
		credential, err := base64.URLEncoding.DecodeString(authhdrs[1])
		if err != nil {
			http.NotFound(w, r)
			return
		}
		credentials := strings.Split(string(credential), ":")
		if len(credentials) != 2 {
			http.NotFound(w, r)
			return
		}
		apikey, err := gouuid.ParseString(credentials[0])
		if err != nil {
			http.NotFound(w, r)
			return
		}
		user, err := umgr.FindByAPIKey(&apikey)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		context.Set(r, "uid", user.UID)
	})
}
