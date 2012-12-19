package main

import (
	"github.com/surma-dump/gouuid"
)

type DomainManager interface {
	ClaimDomain(domain string, uid *gouuid.UUID) error
	FindDomain(domain string) (Domain, error)
	DomainsByUser(uid *gouuid.UUID) ([]Domain, error)
	AddAlias(domain string, alias *Alias) error
}
