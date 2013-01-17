package main

import (
	"github.com/surma-dump/gouuid"
)

// Deprecated to be removed
type Aliases map[string]Alias

type Domain struct {
	Name    string         `bson:"name" json:"name"`
	Owners  []*gouuid.UUID `bson:"owners" json:"-"`
	Aliases []*Alias       `bson:"aliases" json:"aliases"`
}

type Alias struct {
	ID *gouuid.UUID `bson:"id" json:"id"`
	// URL of the repository to link to
	RepoURL string `bson:"repo_url" json:"repo_url"`
	// VCS of the repository ("git", "hg" or "bzr")
	RepoType string `bson:"repo_type" json:"repo_type"`
	// URL to forward to if the URL is being accessed
	// by something else than the go tool.
	ForwardURL string `bson:"forward_url" json:"forward_url"`
	Alias      string `bson:"alias" json:"alias"`
}

type User struct {
	UID            *gouuid.UUID      `bson:"uid" json:"uid"`
	APIKey         *gouuid.UUID      `bson:"apikey" json:"apikey"`
	Authenticators map[string]string `bson:"authenticators" json:"-"`
}
