package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

type Metapage struct {
	domainmgr DomainManager
}

func (m *Metapage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if idx := strings.Index(r.Host, ":"); idx != -1 {
		r.Host = r.Host[:idx]
	}

	if r.URL.Path == "/" {
		domain, err := m.domainmgr.FindDomain(r.Host)
		if err != nil {
			log.Printf("Unknown domain: %s", r.Host)
			http.Redirect(w, r, "http://"+options.Hostname+"/unknown", http.StatusMovedPermanently)
			return
		}
		META_TEMPLATE.Execute(w, map[string]interface{}{
			"Aliases": domain.Aliases,
			"Domain":  r.Host,
		})
		return
	}

	alias, err := m.domainmgr.FindAlias(r.Host, r.URL.Path)
	if err != nil {
		log.Printf("Unknown alias %s|%s: %s", r.Host, r.URL.Path, err)
		http.Redirect(w, r, "http://"+options.Hostname+"/#/unknown", http.StatusMovedPermanently)
		return
	}

	if !isGoGetRequest(r) {
		http.Redirect(w, r, alias.ForwardURL, http.StatusMovedPermanently)
		return
	}
	err = META_TEMPLATE.Execute(w, map[string]interface{}{
		"Aliases": []*Alias{alias},
		"Domain":  r.Host,
	})
	if err != nil {
		log.Printf("Template rendering failed: %s", err)
	}
}

func isGoGetRequest(r *http.Request) bool {
	err := r.ParseForm()
	if err != nil {
		return false
	}
	v, ok := r.Form["go-get"]
	if !ok || len(v) <= 0 {
		return false
	}
	return true
}

var (
	META_TEMPLATE = template.Must(template.New("").Parse(`
		<head>
			{{with $ := .}}
				{{range .Aliases}}
					<meta name="go-import" content="{{$.Domain}}{{.Alias}} {{.RepoType}} {{.RepoURL}}" />
				{{end}}
			{{end}}
		</head>`))
)
