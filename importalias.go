package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/context"

	"github.com/surma-dump/mux"
	"github.com/voxelbrain/goptions"

	"code.google.com/p/goauth2/oauth"

	"labix.org/v2/mgo"
)

var (
	options = struct {
		MongoDB       *url.URL      `goptions:"-m, --mongodb, description='MongoDB to connect to'"`
		ListenAddress *net.TCPAddr  `goptions:"-l, --listen, description='Address to listen on'"`
		Hostname      string        `goptions:"-n, --hostname, obligatory, description='Hostname to serve app on'"`
		StaticDir     string        `goptions:"--static-dir, description='Path to the static content directory'"`
		AuthKeys      []string      `goptions:"--auth-key, description='Add key to an authenticator (format: <authentication provider>:<clientid>:<secret>)'"`
		AuthConfigs   *AuthList      `goptions:"--auth-config, description='Config file for auth apps'"`
		SessionStore  *SessionStore `goptions:"--cookie-key, obligatory, description='Encryption key for cookies'"`
		SessionTTL    time.Duration `goptions:"--session-ttl, description='Duration of a session cookie'"`
		Help          goptions.Help `goptions:"-h, --help, description='Show this help'"`
	}{ // Default values
		MongoDB:       URLMust(url.Parse("mongodb://localhost")),
		ListenAddress: TCPAddrMust(net.ResolveTCPAddr("tcp4", "localhost:8080")),
		StaticDir:     "./static",
		SessionTTL:    30 * time.Minute,
	}
)

func main() {
	goptions.ParseAndFail(&options)

	log.Printf("Connecting to mongodb on %s...", options.MongoDB)
	session, err := mgo.Dial(options.MongoDB.String())
	if err != nil {
		log.Fatalf("Could not connect to %s: %s", options.MongoDB, err)
	}
	defer session.Close()
	db := session.DB("") // Use database specified in URL
	usermgr := &MongoUserManager{db.C("users")}

	mainrouter := mux.NewRouter()
	mainrouter.KeepContext = true
	approuter := mainrouter.Host(options.Hostname).Subrouter()
	api1router := approuter.PathPrefix("/api/v1").Subrouter().StrictSlash(true)

	setupAuthApps(approuter.PathPrefix("/auth").Subrouter(), usermgr)

	api1router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/api/v1")
	})

	approuter.PathPrefix("/").Handler(http.FileServer(http.Dir(options.StaticDir)))
	mainrouter.PathPrefix("/").HandlerFunc(foreignHostname)
	log.Printf("Running webserver...")
	log.Fatalf("Failed to run webserver: %s",
		http.ListenAndServe(options.ListenAddress.String(), mainrouter))
}

type AuthConfig struct {
	Type      string `json:"type"`
	ClientID  string `json:"client_id"`
	Secret    string `json:"secret"`
	AuthURL   string `json:"auth_url"`
	TokenURL  string `json:"token_url"`
	Scope     string `json:"scope"`
	Extractor struct {
		Type  string `json:"type"`
		URL   string `json:"url"`
		Field string `json:"field"`
	} `json:"extractor"`
}

func setupAuthApps(authrouter *mux.Router, usermgr UserManager) {
	enabled := make(AuthList)

	for _, key := range options.AuthKeys {
		keyparts := strings.Split(key, ":")
		if len(keyparts) < 3 {
			log.Printf("Invalid auth key \"%s\" encountered, skipping", key)
			continue
		}
		name, clientid, secret := keyparts[0], keyparts[1], keyparts[2]
		if authconfig, ok := (*options.AuthConfigs)[name]; !ok {
			log.Printf("Unknown authentication provider \"%s\", skipping", keyparts[0])
		} else {
			enabled[name] = authconfig
			enabled[name].ClientID = clientid
			enabled[name].Secret = secret
		}
	}

	for name, authconfig := range enabled {
		var auth Authenticator
		var ex Extractor
		prefix, _ := authrouter.Path("/" + name).URL()
		switch authconfig.Extractor.Type {
		case "json":
			ex = NewJSONExtractor(authconfig.Extractor.URL, authconfig.Extractor.Field)
		default:
			log.Printf("Unknown extractor \"%s\", skipping", authconfig.Extractor.Type)
			continue
		}
		switch authconfig.Type {
		case "oauth":
			log.Printf("Enabling %s OAuth on %s with ClientID %s", name, prefix.String(), authconfig.ClientID)
			auth = NewOAuthAuthenticator(name, &oauth.Config{
				ClientId:     authconfig.ClientID,
				ClientSecret: authconfig.Secret,
				AuthURL:      authconfig.AuthURL,
				TokenURL:     authconfig.TokenURL,
				Scope:        authconfig.Scope,
				RedirectURL:  prefix.String() + "/callback",
			}, ex, usermgr)
		default:
			log.Printf("Unknown authenticator \"%s\", skipping", authconfig.Type)
			continue
		}
		authrouter.PathPrefix("/" + name).Handler(
			context.ClearHandler(HandlerList{
				SilentHandler(SessionOpener(options.SessionStore, int(options.SessionTTL/time.Second))),
				http.StripPrefix(prefix.Path, auth),
				SilentHandler(SessionSaver()),
			}))
	}
}

func TCPAddrMust(t *net.TCPAddr, err error) *net.TCPAddr {
	if err != nil {
		panic(err)
	}
	return t
}

func URLMust(u *url.URL, err error) *url.URL {
	if err != nil {
		panic(err)
	}
	return u
}
