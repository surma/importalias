package main

import (
	"reflect"
	"testing"

	"github.com/surma-dump/gouuid"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

const (
	TOKEN = "s3cr3t"
)

func setup() *mgo.Collection {
	session, err := mgo.Dial("mongodb://localhost/importalias-test")
	if err != nil {
		panic(err)
	}

	db := session.DB("")
	db.DropDatabase()
	db = session.DB("")
	return db.C("user")
}

func teardown(c *mgo.Collection) {
	defer c.Database.Session.Close()
	c.Database.DropDatabase()
}

func TestCreate(t *testing.T) {
	c := setup()
	defer teardown(c)

	mgr := &MongoUserManager{c}
	user, err := mgr.New("gotest", TOKEN)
	if err != nil {
		t.Fatalf("NewUser failed: %s", err)
	}

	qry := c.Find(bson.M{})
	if count, _ := qry.Count(); count != 1 {
		t.Fatalf("Unexpected number of results: %d", count)
	}

	user2 := &User{}
	err = qry.One(user2)
	if err != nil {
		t.Fatalf("Could not get user: %s", err)
	}
	if !reflect.DeepEqual(user, user2) {
		t.Fatalf("Unexpected user data: %#v vs %#v", user, user2)
	}
}

func TestFindUser(t *testing.T) {
	c := setup()
	defer teardown(c)

	mgr := &MongoUserManager{c}
	uuid := gouuid.New()
	err := c.Insert(
		&User{
			APIKey: &uuid,
			Authenticators: map[string]string{
				"gotest": TOKEN,
			},
		})
	if err != nil {
		panic(err)
	}

	user, err := mgr.FindByAuthenticator("gotest", TOKEN)
	if err != nil {
		t.Fatalf("Could not get user: %s", err)
	}
	if !reflect.DeepEqual(user.APIKey, &uuid) {
		t.Fatalf("Unexpected APIKey for user")
	}

	user2, err := mgr.FindByAPIKey(&uuid)
	if err != nil {
		t.Fatalf("Could not get user: %s", err)
	}
	if !reflect.DeepEqual(user, user2) {
		t.Fatalf("Unexpected user data")
	}
}

func TestAddAuthenticator(t *testing.T) {
	c := setup()
	defer teardown(c)

	mgr := &MongoUserManager{c}
	user, err := mgr.New("gotest", TOKEN)
	if err != nil {
		panic(err)
	}

	err = mgr.AddAuthenticator(user.UID, "gotest2", TOKEN)
	if err != nil {
		t.Fatalf("Could not add second authenticator: %s", err)
	}

	user2, err := mgr.FindByAuthenticator("gotest2", TOKEN)
	if err != nil {
		t.Fatalf("Could not get user: %s", err)
	}
	if !reflect.DeepEqual(user.UID, user2.UID) {
		t.Fatalf("Unexpected APIKey for user")
	}
}
