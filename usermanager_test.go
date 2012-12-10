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

	mgr := NewMongoUserManager(c)
	err := mgr.New("gotest", TOKEN)
	if err != nil {
		t.Fatalf("NewUser failed: %s", err)
	}

	qry := c.Find(bson.M{})
	if count, _ := qry.Count(); count != 1 {
		t.Fatalf("Unexpected number of results: %d", count)
	}

	user := &User{}
	err = qry.One(user)
	if err != nil {
		t.Fatalf("Could not get user: %s", err)
	}
	if !(user.APIKey.String() != "" &&
		user.Authenticators["gotest"] == TOKEN) {
		t.Fatalf("Unexpected user data: %#v", user)
	}
}

func TestFindUser(t *testing.T) {
	c := setup()
	defer teardown(c)

	mgr := NewMongoUserManager(c)
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

	mgr := NewMongoUserManager(c)
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

	err = mgr.AddAuthenticator(&uuid, "gotest2", TOKEN)
	if err != nil {
		t.Fatalf("Could not add second authenticator: %s", err)
	}

	user, err := mgr.FindByAuthenticator("gotest2", TOKEN)
	if err != nil {
		t.Fatalf("Could not get user: %s", err)
	}
	if !reflect.DeepEqual(user.APIKey, &uuid) {
		t.Fatalf("Unexpected APIKey for user")
	}
}
