package main

import (
	"github.com/surma-dump/gouuid"
)

type Domain struct {
	Name    string         `bson:"name" json:"name"`
	Owners  []*gouuid.UUID `bson:"owners" json:"-"`
	Aliases []*Alias       `bson:"aliases" json:"aliases"`
}

func (d *Domain) IsOwnedBy(uid *gouuid.UUID) bool {
	for _, owner := range d.Owners {
		if owner.Equal(*uid) {
			return true
		}
	}
	return false
}
