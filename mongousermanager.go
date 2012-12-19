package main

import (
	"fmt"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type MongoUserManager struct {
	Collection *mgo.Collection
}

func (mum *MongoUserManager) FindByAuthenticator(authenticator, id string) (*User, error) {
	user := &User{}
	qry := mum.Collection.Find(bson.M{
		"authenticators." + authenticator: id,
	})
	if count, _ := qry.Count(); count != 1 {
		return nil, ErrNotFound
	}
	err := qry.One(user)
	return user, err
}

func (mum *MongoUserManager) FindByAPIKey(apikey *gouuid.UUID) (*User, error) {
	user := &User{}
	qry := mum.Collection.Find(bson.M{
		"apikey": apikey,
	})
	if count, _ := qry.Count(); count != 1 {
		return nil, ErrNotFound
	}
	err := qry.One(user)
	return user, err
}

func (mum *MongoUserManager) FindByUID(uid *gouuid.UUID) (*User, error) {
	user := &User{}
	qry := mum.Collection.Find(bson.M{
		"uid": uid,
	})
	if count, _ := qry.Count(); count != 1 {
		return nil, ErrNotFound
	}
	err := qry.One(user)
	return user, err
}

func (mum *MongoUserManager) New(authenticator, id string) (*User, error) {
	uid, apikey := gouuid.New(), gouuid.New()
	user := &User{
		UID:    &uid,
		APIKey: &apikey,
		Authenticators: map[string]string{
			authenticator: id,
		},
	}
	return user, mum.UpdateUser(user)
}

func (mum *MongoUserManager) UpdateUser(u *User) error {
	_, err := mum.Collection.Upsert(bson.M{
		"uid": u.UID,
	}, u)
	return err
}

func (mum *MongoUserManager) AddAuthenticator(uid *gouuid.UUID, authenticator, id string) error {
	return mum.Collection.Update(bson.M{
		"uid": uid,
	}, bson.M{
		"$set": bson.M{
			"authenticators": bson.M{
				authenticator: id,
			},
		},
	})
}
