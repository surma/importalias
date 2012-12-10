package main

import (
	"fmt"

	"github.com/surma-dump/gouuid"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

var (
	ErrNotFound = fmt.Errorf("User not found")
)

type UserManager interface {
	FindByAuthenticator(authenticator, id string) (*User, error)
	FindByAPIKey(apikey *gouuid.UUID) (*User, error)
	New(authenticator, id string) error
	AddAuthenticator(apikey *gouuid.UUID, authenticator, id string) error
}

type MongoUserManager struct {
	collection *mgo.Collection
}

func NewMongoUserManager(c *mgo.Collection) *MongoUserManager {
	return &MongoUserManager{c}
}

func (mum *MongoUserManager) FindByAuthenticator(authenticator, id string) (*User, error) {
	user := &User{}
	qry := mum.collection.Find(bson.M{
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
	qry := mum.collection.Find(bson.M{
		"apikey": apikey,
	})
	if count, _ := qry.Count(); count != 1 {
		return nil, ErrNotFound
	}
	err := qry.One(user)
	return user, err
}

func (mum *MongoUserManager) New(authenticator, id string) error {
	uuid := gouuid.New()
	return mum.UpdateUser(&User{
		APIKey: &uuid,
		Authenticators: map[string]string{
			authenticator: id,
		},
	})
}

func (mum *MongoUserManager) UpdateUser(u *User) error {
	_, err := mum.collection.Upsert(bson.M{
		"apikey": u.APIKey,
	}, u)
	return err
}

func (mum *MongoUserManager) AddAuthenticator(apikey *gouuid.UUID, authenticator, id string) error {
	return mum.collection.Update(bson.M{
		"apikey": apikey,
	}, bson.M{
		"$set": bson.M{
			"authenticators": bson.M{
				authenticator: id,
			},
		},
	})
}
