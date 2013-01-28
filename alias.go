package main

import (
	"github.com/surma-dump/gouuid"
)

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
