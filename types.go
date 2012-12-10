package main

import (
	"github.com/surma-dump/gouuid"
)

type Aliases map[string]Alias
type Alias struct {
	// User ids which can access this domain
	Owners []string
	// URL of the repository to link to
	RepoURL string `json:"repo_url"`
	// VCS of the repository ("git", "hg" or "bzr")
	RepoType string `json:"repo_type"`
	// URL to forward to if the URL is being accessed
	// by something else than the go tool.
	ForwardURL string `josn:"forward_url"`
	Alias      string `json:"alias"`
}

type User struct {
	APIKey         *gouuid.UUID      `json:"apikey"`
	Authenticators map[string]string `json:"authenticators"`
}
