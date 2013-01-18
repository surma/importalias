package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/surma-dump/context"
	"github.com/surma-dump/gouuid"

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
	api.Router.Path("/domains/{domain}").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.ListAliases(w, r)
	})
	api.Router.Path("/domains/{domain}").Methods("DELETE").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.DeleteDomain(w, r)
	})
	api.Router.Path("/domains/{domain}").Methods("PUT").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.SetAlias(w, r)
	})
	api.Router.Path("/domains/{domain}/{aid}").Methods("DELETE").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.DeleteAlias(w, r)
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
	m := new(dns.Msg)
	m.SetQuestion(domain, dns.TypeTXT)
	c := new(dns.Client)
	rr, _, err := c.Exchange(m, "8.8.8.8:53")
	if err != nil {
		log.Printf("Could not check domain’s TXT record: %s", err)
		http.Error(w, "Could not check domain’s TXT record", http.StatusInternalServerError)
		return
	}

	uid := context.Get(r, "uid").(*gouuid.UUID)
	if !containsValidTXTRecord(uid, rr.Answer) {
		http.Error(w, "Did not find your UID in domain’s TXT records", http.StatusUnauthorized)
		return
	}

	err = api.domainmgr.ClaimDomain(vars["domain"], uid)
	if err == ErrAlreadyClaimed {
		http.Error(w, "Domain already claimed", http.StatusForbidden)
		return
	} else if err != nil {
		log.Printf("Could not claim domain: %s", err)
		http.Error(w, "Could not claim domain", http.StatusInternalServerError)
		return
	}
	http.Error(w, "", http.StatusNoContent)
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

func (api *APIv1) DeleteDomain(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid := context.Get(r, "uid").(*gouuid.UUID)
	err := api.domainmgr.DeleteDomain(vars["domain"], uid)
	if err != nil {
		http.Error(w, "Could not delete domain", http.StatusInternalServerError)
		return
	}
	http.Error(w, "", http.StatusNoContent)
}

func (api *APIv1) ListAliases(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid := context.Get(r, "uid").(*gouuid.UUID)

	domain, err := api.domainmgr.FindDomain(vars["domain"])
	if err != nil || !domain.IsOwnedBy(uid) {
		http.Error(w, "Could not find domain", http.StatusNotFound)
		return
	}
	enc := json.NewEncoder(w)
	enc.Encode(domain.Aliases)
}

func (api *APIv1) SetAlias(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid := context.Get(r, "uid").(*gouuid.UUID)

	alias := &Alias{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(alias)
	if err != nil {
		http.Error(w, "Invalid alias object", http.StatusBadRequest)
		return
	}

	err = api.domainmgr.SetAlias(vars["domain"], alias, uid)
	if err != nil {
		http.Error(w, "Could not add alias", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusCreated)
	enc := json.NewEncoder(w)
	enc.Encode(alias.ID)
}

func (api *APIv1) DeleteAlias(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid := context.Get(r, "uid").(*gouuid.UUID)

	aid, err := gouuid.ParseString(vars["aid"])
	if err != nil {
		http.Error(w, "Invalid id format", http.StatusNotFound)
		return
	}

	err = api.domainmgr.DeleteAlias(&aid, uid)
	if err != nil {
		http.Error(w, "Could not delete alias", http.StatusNotFound)
		return
	}
	http.Error(w, "", http.StatusNoContent)
}
