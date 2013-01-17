package main

import (
	"fmt"

	"github.com/surma-dump/gouuid"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type MongoDomainManager struct {
	Collection *mgo.Collection
}

func (mdm *MongoDomainManager) ClaimDomain(name string, uid *gouuid.UUID) error {
	_, err := mdm.FindDomain(name)
	if err == nil {
		return ErrAlreadyClaimed
	} else if err != ErrNotFound {
		return fmt.Errorf("Invalid domain query result")
	}

	return mdm.Collection.Insert(&Domain{
		Name:    name,
		Owners:  []*gouuid.UUID{uid},
		Aliases: []*Alias{},
	})
}

func (mdm *MongoDomainManager) FindDomain(name string) (*Domain, error) {
	domain := &Domain{}
	qry := mdm.Collection.Find(bson.M{
		"name": name,
	})
	if count, _ := qry.Count(); count != 1 {
		return nil, ErrNotFound
	}
	err := qry.One(domain)
	return domain, err
}

func (mdm *MongoDomainManager) DomainsByOwner(uid *gouuid.UUID) ([]*Domain, error) {
	domains := []*Domain{}
	qry := mdm.Collection.Find(bson.M{
		"owners": uid,
	})
	// FIXME: This can crash everything if there are *a lot* of
	// domains. Paging?
	err := qry.All(&domains)
	return domains, err
}

func (mdm *MongoDomainManager) SetAlias(name string, alias *Alias) error {
	domain, err := mdm.FindDomain(name)
	if err != nil {
		return err
	}
	aid := gouuid.New()
	alias.ID = &aid
	_, err = mdm.Collection.Upsert(bson.M{
		"name": domain.Name,
	}, bson.M{
		"$push": bson.M{
			"aliases": alias,
		},
	})
	return err
}

func (mdm *MongoDomainManager) DeleteAlias(aid *gouuid.UUID) error {
	_, err := mdm.Collection.Upsert(bson.M{
		"aliases.id": aid,
	}, bson.M{
		"$pull": bson.M{
			"aliases": bson.M{
				"id": aid,
			},
		},
	})
	return err
}
