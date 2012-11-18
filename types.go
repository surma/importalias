package main

type Aliases map[string]Alias
type Alias struct {
	RepoURL    string `json:"repo_url"`
	RepoType   string `json:"repo_type"`
	ForwardURL string `josn:"forward_url"`
	Alias      string `json:"alias"`
}
