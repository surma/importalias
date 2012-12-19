package main

import (
	"github.com/surma-dump/gouuid"
)

// Deprecated to be removed
type Aliases map[string]Alias

type Domain struct {
	Name    string           `json:"name"`
	Owners  []string         `bson:"owners"`
	Aliases map[string]Alias `json:"aliases"`
}

type Alias struct {
	ID *gouuid.UUID `json:"id"`
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
	UID            *gouuid.UUID      `json:"uid"`
	APIKey         *gouuid.UUID      `json:"apikey"`
	Authenticators map[string]string `json:"authenticators"`
}
