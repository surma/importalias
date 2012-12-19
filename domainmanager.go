package main

import (
	"github.com/surma-dump/gouuid"
)

type DomainManager interface {
	ClaimDomain(domain string, uid *gouuid.UUID) error
	FindDomain(domain string) (*Domain, error)
	DomainsByOwner(uid *gouuid.UUID) ([]*Domain, error)
	SetAlias(domain string, alias *Alias) error
	DeleteAlias(aid *gouuid.UUID) error
}
