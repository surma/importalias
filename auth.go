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

	"github.com/gorilla/sessions"
	"github.com/surma-dump/context"
	"github.com/surma-dump/gouuid"
)

func init() {
	// Register *gouuid.UUID as a type with gob
	// as gob is being used bei gorilla’s cookie session store which
	// is used by us to save the user’s uuid.
	uuid := gouuid.New()
	gob.Register(&uuid)
}

var (
	CALLBACK_TEMPLATE = template.Must(template.New("").Parse(`
		<script>
			if(window.opener) {
				window.opener.postMessage("auth_done", "*");
				window.close();
			}
		</script>`))
)

func LoginHandler(umgr UserManager) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer CALLBACK_TEMPLATE.Execute(w, nil)
		authname, ok1 := context.Get(r, "authname").(string)
		uid, ok2 := context.Get(r, "authuid").(string)
		if !(ok1 && ok2) {
			return
		}
		// Login
		user, err := umgr.FindByAuthenticator(authname, uid)
		if err != nil && err != ErrUserNotFound {
			user = nil
			log.Printf("Could not query user database: %s", err)
			return
		}
		if err == ErrUserNotFound {
			// New user
			user, err = umgr.New(authname, uid)
			if err != nil {
				log.Printf("Error creating user: %s", err)
				return
			}
		}
		session := context.Get(r, "session").(*sessions.Session)
		session.Values["uid"] = user.UID
		session.Save(r, w)
	})
}

func LogoutHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer CALLBACK_TEMPLATE.Execute(w, nil)
		session := context.Get(r, "session").(*sessions.Session)
		session.Values = make(map[interface{}]interface{})
		session.Save(r, w)

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

func authListHandler(auths *AuthList) http.Handler {
	authnames := make([]string, 0, len(*auths))
	for name, authconfig := range *auths {
		if authconfig.AuthKey != nil {
			authnames = append(authnames, name)
		}
	}
	authlist, err := json.Marshal(authnames)
	if err != nil {
		panic("Could not create auth list handler: " + err.Error())
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(authlist)
	})
}
