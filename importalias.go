package main

import (
	"github.com/gorilla/mux"
	"github.com/voxelbrain/goptions"
	"labix.org/v2/mgo"
	"log"
	"net"
	"net/http"
	"net/url"
)

var (
	options = struct {
		MongoDB       *url.URL     `goptions:"-m, --mongodb, description='MongoDB to connect to'"`
		ListenAddress *net.TCPAddr `goptions:"-l, --listen, description='Address to listen on'"`
		Hostname      string       `goptions:"-n, --hostname, obligatory, description='Hostname to serve app on'"`
		StaticDir     string       `goptions:"--static-dir, description='Path to the static content directory'"`
		goptions.Help `goptions:"-h, --help, description='Show this help'"`
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
	approuter.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))
	apirouter := approuter.PathPrefix("/api").Subrouter()
	_, _ = db, apirouter
	log.Printf("Running webserver...")
	log.Fatalf("Failed to run webserver: %s", http.ListenAndServe(options.ListenAddress.String(), mainrouter))
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
