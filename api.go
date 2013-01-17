package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/context"
	"github.com/surma-dump/gouuid"
	"github.com/surma-dump/mux"

	"github.com/miekg/dns"
)

type API interface {
	http.Handler
}

type APIv1 struct {
	domainmgr DomainManager
	usermgr   UserManager
	*mux.Router
}

func NewAPIv1(domainmgr DomainManager, usermgr UserManager) *APIv1 {
	api := &APIv1{
		domainmgr: domainmgr,
		usermgr:   usermgr,
	}

	api.Router = mux.NewRouter()
	api.Router.KeepContext = true
	api.Router.Path("/me").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uid := context.Get(r, "uid").(*gouuid.UUID)
		user, err := usermgr.FindByUID(uid)
		if err != nil {
			log.Printf("Session UID invalid: %s", err)
			http.Error(w, "Could not find user", http.StatusInternalServerError)
			return
		}
		enc := json.NewEncoder(w)
		enc.Encode(user)
	})

	api.Router.Path("/domains").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.GetDomains(w, r)
	})
	api.Router.Path("/domains/{domain}").Methods("POST").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.ClaimDomain(w, r)
	})

	return api
}

func (api *APIv1) GetDomains(w http.ResponseWriter, r *http.Request) {
	uid := context.Get(r, "uid").(*gouuid.UUID)
	domains, err := api.domainmgr.DomainsByOwner(uid)
	if err != nil {
		log.Printf("Could not list domains: %s", err)
		http.Error(w, "Could not list domains", http.StatusInternalServerError)
		return
	}
	domainnames := make([]string, len(domains))
	for i, domain := range domains {
		domainnames[i] = domain.Name
	}

	enc := json.NewEncoder(w)
	enc.Encode(domainnames)
}

func (api *APIv1) ClaimDomain(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	domain := dns.Fqdn("_importalias." + vars["domain"])
	rr, err := RecursiveDNS(domain, dns.TypeTXT)
	if err != nil {
		log.Printf("Could not check domain’s TXT record: %s", err)
		http.Error(w, "Could not check domain’s TXT record", http.StatusInternalServerError)
		return
	}

	uid := context.Get(r, "uid").(*gouuid.UUID)
	if !containsValidTXTRecord(uid, rr) {
		http.Error(w, "Did not find your UID in domain’s TXT records", http.StatusUnauthorized)
		return
	}

	log.Printf("Claimed! (not really)")
	return
}

func containsValidTXTRecord(uid *gouuid.UUID, rr []dns.RR) bool {
	for _, record := range rr {
		if txtr, ok := record.(*dns.TXT); ok {
			for _, txt := range txtr.Txt {
				txtuid, err := gouuid.ParseString(txt)
				if err != nil {
					continue
				}
				if txtuid.Equal(*uid) {
					return true
				}
			}
		}
	}
	return false
}
