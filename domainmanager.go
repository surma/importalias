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
	DeleteDomain(domain string, uid *gouuid.UUID) error
	DomainsByOwner(uid *gouuid.UUID) ([]*Domain, error)
	FindAlias(domain string, alias string) (*Alias, error)
	SetAlias(domain string, alias *Alias, uid *gouuid.UUID) error
	DeleteAlias(aid *gouuid.UUID) error
}
