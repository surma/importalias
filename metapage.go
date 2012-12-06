package main

import (
	"html/template"
	"net/http"
	"strings"
)

var (
	knownHosts = map[string]Aliases{
		"localhostalias.de": Aliases{
			"/goptions": Alias{
				RepoURL:    "https://github.com/surma/goptions",
				ForwardURL: "https://github.com/surma/goptions",
				RepoType:   "git",
				Alias:      "localhostalias.de/goptions",
			},
			"/gocpio": Alias{
				RepoURL:    "https://github.com/surma/gocpio",
				ForwardURL: "https://github.com/surma/gocpio",
				RepoType:   "git",
				Alias:      "localhostalias.de/gocpio",
			},
		},
	}
)

func foreignHostname(w http.ResponseWriter, r *http.Request) {
	if idx := strings.Index(r.Host, ":"); idx != -1 {
		r.Host = r.Host[:idx]
	}
	aliases, ok := knownHosts[r.Host]
	if !ok {
		http.Redirect(w, r, "http://"+options.Hostname+"/unknown", http.StatusMovedPermanently)
		return
	}
	if r.URL.Path == "/" {
		TEMPLATE.Execute(w, aliases)
		return
	}

	alias, ok := aliases[r.URL.Path]
	if !ok {
		http.Redirect(w, r, "http://"+options.Hostname+"/unknown", http.StatusMovedPermanently)
		return
	}

	if !isGoGetRequest(r) {
		http.Redirect(w, r, alias.ForwardURL, http.StatusMovedPermanently)
		return
	}
	TEMPLATE.Execute(w, []Alias{alias})
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
	TEMPLATE = template.Must(template.New("").Parse(`
		<head>
			{{range .}}
			<meta name="go-import" content="{{.Alias}} {{.RepoType}} {{.RepoURL}}" />
			{{end}}
		</head>`))
)
