package main

import (
	"github.com/surma-dump/gouuid"
)

var (
	ErrNotFound = fmt.Errorf("User not found")
)

type UserManager interface {
	FindByAuthenticator(authenticator, id string) (*User, error)
	FindByAPIKey(apikey *gouuid.UUID) (*User, error)
	FindByUID(uid *gouuid.UUID) (*User, error)
	New(authenticator, id string) (*User, error)
	AddAuthenticator(uid *gouuid.UUID, authenticator, id string) error
}
