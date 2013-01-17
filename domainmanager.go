package main

import (
	"fmt"

	"github.com/surma-dump/gouuid"
)

var (
	ErrAlreadyClaimed = fmt.Errorf("Domain already claimed")
)

type DomainManager interface {
	ClaimDomain(domain string, uid *gouuid.UUID) error
	FindDomain(domain string) (*Domain, error)
	DomainsByOwner(uid *gouuid.UUID) ([]*Domain, error)
	SetAlias(domain string, alias *Alias) error
	DeleteAlias(aid *gouuid.UUID) error
}
