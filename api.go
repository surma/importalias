package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/context"
	"github.com/surma-dump/gouuid"
	"github.com/surma-dump/mux"
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
	api.Router.Methods("GET").Path("/domains").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.GetDomains(w, r)
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
