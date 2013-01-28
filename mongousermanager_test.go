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

func umgr_setup() *mgo.Collection {
	session, err := mgo.Dial("mongodb://localhost/importalias-test")
	if err != nil {
		panic(err)
	}

	db := session.DB("")
	db.C("user").DropCollection()
	return db.C("user")
}

func umgr_teardown(c *mgo.Collection, t *testing.T) {
	if t.Failed() {
		var v interface{}
		it := c.Find(bson.M{}).Iter()
		t.Logf("Datbase documents:")
		for it.Next(&v) {
			t.Logf("%#v", v)
			v = nil
		}
		if it.Err() != nil {
			t.Logf("Error: %s", it.Err())
		}
	}
	defer c.Database.Session.Close()
	c.DropCollection()
}

func TestCreate(t *testing.T) {
	c := umgr_setup()
	defer umgr_teardown(c, t)

	mgr := UserManager(&MongoUserManager{c})
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
		t.Fatalf("Unexpected user data: Got %#v, expected %#v", user2, user)
	}
}

func TestFindUser(t *testing.T) {
	c := umgr_setup()
	defer umgr_teardown(c, t)

	mgr := UserManager(&MongoUserManager{c})
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
	c := umgr_setup()
	defer umgr_teardown(c, t)

	mgr := UserManager(&MongoUserManager{c})
	user, err := mgr.New("gotest", TOKEN)
	if err != nil {
		panic(err)
	}

	err = mgr.AddAuthenticator(user.UID, "gotest2", TOKEN)
	if err != nil {
		t.Fatalf("Could not add second authenticator: %s", err)
	}

	user2, err := mgr.FindByAuthenticator("gotest", TOKEN)
	if err != nil {
		t.Fatalf("Could not get user: %s", err)
	}
	if !reflect.DeepEqual(user.UID, user2.UID) {
		t.Fatalf("Could not find user with old authenticator")
	}

	user2, err = mgr.FindByAuthenticator("gotest2", TOKEN)
	if err != nil {
		t.Fatalf("Could not get user: %s", err)
	}
	if !reflect.DeepEqual(user.UID, user2.UID) {
		t.Fatalf("Could not find user with new authenticator")
	}
}
