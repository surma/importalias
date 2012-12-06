package main

type Aliases map[string]Alias
type Alias struct {
	// URL of the repository to link to
	RepoURL string `json:"repo_url"`
	// VCS of the repository ("git", "hg" or "bzr")
	RepoType string `json:"repo_type"`
	// URL to forward to if the URL is being accessed
	// by something else than the go tool.
	ForwardURL string `josn:"forward_url"`
	Alias      string `json:"alias"`
}
