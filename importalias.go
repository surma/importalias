package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"path"

	"github.com/gorilla/mux"
	"github.com/surma-dump/goauth2/oauth"
	"github.com/voxelbrain/goptions"
	"labix.org/v2/mgo"
)

var (
	options = struct {
		MongoDB        *url.URL     `goptions:"-m, --mongodb, description='MongoDB to connect to'"`
		ListenAddress  *net.TCPAddr `goptions:"-l, --listen, description='Address to listen on'"`
		Hostname       string       `goptions:"-n, --hostname, obligatory, description='Hostname to serve app on'"`
		StaticDir      string       `goptions:"--static-dir, description='Path to the static content directory'"`
		GitHubClientID string       `goptions:"--github-clientid, description='Client ID of the GitHub App'"`
		GitHubSecret   string       `goptions:"--github-secret, description='Secret of the GitHub App'"`
		GoogleClientID string       `goptions:"--google-clientid, description='Client ID of the Google App'"`
		GoogleSecret   string       `goptions:"--google-secret, description='Secret of the Google App'"`
		goptions.Help  `goptions:"-h, --help, description='Show this help'"`
	}{ // Default values
		MongoDB:       URLMust(url.Parse("mongodb://localhost")),
		ListenAddress: TCPAddrMust(net.ResolveTCPAddr("tcp4", "localhost:8080")),
		StaticDir:     "./static",
	}
)

func init() {
	goptions.ParseAndFail(&options)
}

func main() {
	session, err := mgo.Dial(options.MongoDB.String())
	if err != nil {
		log.Fatalf("Could not connect to %s: %s", options.MongoDB, err)
	}
	defer session.Close()
	db := session.DB("") // Use database specified in URL

	mainrouter := mux.NewRouter()
	approuter := mainrouter.Host(options.Hostname).Subrouter()
	api1router := approuter.PathPrefix("/api/v1").Subrouter().StrictSlash(true)

	setupAuthApps(approuter)

	api1router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "/api/v1")
	})

	approuter.PathPrefix("/").Handler(http.FileServer(http.Dir(options.StaticDir)))
	mainrouter.PathPrefix("/").HandlerFunc(foreignHostname)
	_ = db
	log.Printf("Running webserver...")
	log.Fatalf("Failed to run webserver: %s",
		http.ListenAndServe(options.ListenAddress.String(), mainrouter))
}

func setupAuthApps(approuter *mux.Router) {
	if len(options.GitHubClientID) > 0 && len(options.GitHubSecret) > 0 {
		log.Printf("Enabling GitHub auth with ClientID %s", options.GitHubClientID)
		approuter.PathPrefix("/auth/github").
			Handler(http.StripPrefix("/auth/github", NewOAuthAuthenticator(&oauth.Config{
			ClientId:     options.GitHubClientID,
			ClientSecret: options.GitHubSecret,
			AuthURL:      "https://github.com/login/oauth/authorize",
			TokenURL:     "https://github.com/login/oauth/access_token",
			RedirectURL:  "http://" + path.Join(options.Hostname, "/auth/github/callback"),
		}, ExtractorFunc(GitHubExtractor))))
	}
	if len(options.GoogleClientID) > 0 && len(options.GoogleSecret) > 0 {
		log.Printf("Enabling Google auth with ClientID %s", options.GoogleClientID)
		approuter.PathPrefix("/auth/google").
			Handler(http.StripPrefix("/auth/google", NewOAuthAuthenticator(&oauth.Config{
			ClientId:     options.GoogleClientID,
			ClientSecret: options.GoogleSecret,
			Scope:        "https://www.googleapis.com/auth/userinfo.email",
			AuthURL:      "https://accounts.google.com/o/oauth2/auth",
			TokenURL:     "https://accounts.google.com/o/oauth2/token",
			RedirectURL:  "http://" + path.Join(options.Hostname, "/auth/google/callback"),
		}, ExtractorFunc(GoogleExtractor))))
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
