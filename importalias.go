package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
	"github.com/surma-dump/context"
	"github.com/voxelbrain/goptions"

	oauth2 "code.google.com/p/goauth2/oauth"

	"labix.org/v2/mgo"
)

var (
	options = struct {
		MongoDB       *url.URL      `goptions:"-m, --mongodb, description='MongoDB to connect to'"`
		ListenAddress *net.TCPAddr  `goptions:"-l, --listen, description='Address to listen on'"`
		Hostname      string        `goptions:"-n, --hostname, obligatory, description='Hostname to serve app on'"`
		StaticDir     string        `goptions:"--static-dir, description='Path to the static content directory'"`
		AuthConfigs   *AuthList     `goptions:"--auth-config, description='Config file for auth apps'"`
		AuthKeys      []*AuthKey    `goptions:"--auth-key, description='Add key to an authenticator (format: <authentication provider>:<clientid>:<secret>)'"`
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
	domainmgr := &MongoDomainManager{db.C("domains")}

	mainrouter := mux.NewRouter()
	approuter := mainrouter.Host(options.Hostname).Subrouter()
	authrouter := approuter.PathPrefix("/auth").Subrouter()
	apirouter := approuter.PathPrefix("/api").Subrouter()

	setupAuthApps(authrouter, usermgr)
	setupApiApps(apirouter, domainmgr, usermgr)

	approuter.PathPrefix("/").Handler(http.FileServer(http.Dir(options.StaticDir)))
	mainrouter.PathPrefix("/").Handler(&Metapage{domainmgr})
	log.Printf("Running webserver...")
	log.Fatalf("Failed to run webserver: %s", http.ListenAndServe(options.ListenAddress.String(), mainrouter))
}

func setupAuthApps(authrouter *mux.Router, usermgr UserManager) {
	for _, authkey := range options.AuthKeys {
		authconfig, ok := (*options.AuthConfigs)[authkey.Name]
		if !ok {
			log.Printf("Unknown authenticator \"%s\", skipping", authkey.Name)
			continue
		}
		prefix, _ := authrouter.Path("/" + authkey.Name).URL()
		authconfig.AuthKey = authkey
		authconfig.RedirectURL = prefix.String() + "/callback"

		if authapp := createAuthApp(authconfig); authapp != nil {
			log.Printf("Mounting %s authentication service on %s with ClientID %s", authconfig.AuthKey.Name, prefix.String(), authconfig.AuthKey.ClientID)
			authrouter.PathPrefix("/" + authkey.Name).Handler(
				context.ClearHandler(HandlerList{
					SilentHandler(SessionHandler(options.SessionStore, int(options.SessionTTL/time.Second))),
					SilentHandler(http.StripPrefix(prefix.Path, authapp)),
					LoginHandler(usermgr),
				}))
		}
	}
	authrouter.Handle("/", authListHandler(options.AuthConfigs))
	authrouter.Path("/logout").Handler(
		context.ClearHandler(HandlerList{
			SilentHandler(SessionHandler(options.SessionStore, int(options.SessionTTL/time.Second))),
			LogoutHandler(),
		}))
}

func createAuthApp(authconfig *AuthConfig) (auth AuthenticationService) {
	var ex Extractor
	switch authconfig.Extractor.Type {
	case "json":
		ex = NewJSONExtractor(authconfig.Extractor.URL, authconfig.Extractor.Field)
	default:
		log.Printf("Unknown extractor \"%s\", skipping", authconfig.Extractor.Type)
		return
	}
	switch authconfig.Type {
	case "oauth2":
		auth = NewOAuth2AuthenticationService(authconfig.AuthKey.Name, &oauth2.Config{
			ClientId:     authconfig.AuthKey.ClientID,
			ClientSecret: authconfig.AuthKey.Secret,
			AuthURL:      authconfig.AuthURL,
			TokenURL:     authconfig.TokenURL,
			Scope:        authconfig.Scope,
			RedirectURL:  authconfig.RedirectURL,
		}, ex)
	default:
		panic(fmt.Sprintf("Unknown authenticator \"%s\", skipping", authconfig.Type))
	}
	return
}

func setupApiApps(apirouter *mux.Router, domainmgr DomainManager, usermgr UserManager) {
	prefix, _ := apirouter.Path("/").URL()
	apirouter.PathPrefix("/v1").Handler(
		context.ClearHandler(HandlerList{
			SilentHandler(SessionHandler(options.SessionStore, int(options.SessionTTL/time.Second))),
			SilentHandler(BasicAuth(usermgr)),
			SilentHandler(ValidateUID(usermgr)),
			http.StripPrefix(prefix.Path+"v1", NewAPIv1(domainmgr, usermgr)),
		}))
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
