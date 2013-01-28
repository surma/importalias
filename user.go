package main

import (
	"github.com/surma-dump/gouuid"
)

type User struct {
	UID            *gouuid.UUID      `bson:"uid" json:"uid"`
	APIKey         *gouuid.UUID      `bson:"apikey" json:"apikey"`
	Authenticators map[string]string `bson:"authenticators" json:"-"`
}
