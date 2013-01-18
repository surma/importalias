package main

import (
	"reflect"
	"testing"

	"github.com/surma-dump/gouuid"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

const ()

func dmgr_setup() *mgo.Collection {
	session, err := mgo.Dial("mongodb://localhost/importalias-test")
	if err != nil {
		panic(err)
	}

	db := session.DB("")
	db.DropDatabase()
	db = session.DB("")
	return db.C("domains")
}

func dmgr_teardown(c *mgo.Collection, t *testing.T) {
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
	c.Database.DropDatabase()
}

func TestClaim(t *testing.T) {
	c := dmgr_setup()
	defer dmgr_teardown(c, t)

	mgr := DomainManager(&MongoDomainManager{c})
	uid := gouuid.New()
	err := mgr.ClaimDomain("gotest.org", &uid)
	if err != nil {
		t.Fatalf("ClaimDomain failed: %s", err)
	}

	qry := c.Find(bson.M{})
	if count, _ := qry.Count(); count != 1 {
		t.Fatalf("Unexpected number of results: %d", count)
	}

	got, expected := Domain{}, Domain{
		Name:    "gotest.org",
		Owners:  []*gouuid.UUID{&uid},
		Aliases: []*Alias{},
	}
	err = qry.One(&got)
	if err != nil {
		t.Fatalf("Could not get domain: %s", err)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Unexpected domain data. Got %#v, expected %#v", got, expected)
	}
}

func TestMultiClaim(t *testing.T) {
	c := dmgr_setup()
	defer dmgr_teardown(c, t)

	mgr := DomainManager(&MongoDomainManager{c})
	uid := gouuid.New()
	err := mgr.ClaimDomain("gotest1.org", &uid)
	if err != nil {
		t.Fatalf("ClaimDomain failed: %s", err)
	}
	err = mgr.ClaimDomain("gotest2.org", &uid)
	if err != nil {
		t.Fatalf("ClaimDomain failed: %s", err)
	}
	err = mgr.ClaimDomain("gotest3.org", &uid)
	if err != nil {
		t.Fatalf("ClaimDomain failed: %s", err)
	}

	qry := c.Find(bson.M{})
	if count, _ := qry.Count(); count != 3 {
		t.Fatalf("Unexpected number of results: %d", count)
	}

	got, expected := []Domain{}, []Domain{
		Domain{
			Name:    "gotest1.org",
			Owners:  []*gouuid.UUID{&uid},
			Aliases: []*Alias{},
		},
		Domain{
			Name:    "gotest2.org",
			Owners:  []*gouuid.UUID{&uid},
			Aliases: []*Alias{},
		},
		Domain{
			Name:    "gotest3.org",
			Owners:  []*gouuid.UUID{&uid},
			Aliases: []*Alias{},
		},
	}
	err = qry.All(&got)
	if err != nil {
		t.Fatalf("Could not get domain: %s", err)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Unexpected domain data. Got %#v, expected %#v", got, expected)
	}
}

func TestFindDomain(t *testing.T) {
	c := dmgr_setup()
	defer dmgr_teardown(c, t)

	mgr := DomainManager(&MongoDomainManager{c})
	uid, other_uid := gouuid.New(), gouuid.New()
	mgr.ClaimDomain("gotest.org", &uid)
	mgr.ClaimDomain("gotest2.org", &other_uid)

	got, err := mgr.FindDomain("gotest.org")
	if err != nil {
		t.Fatalf("Could find domain: %s", err)
	}

	expected := &Domain{
		Name:    "gotest.org",
		Owners:  []*gouuid.UUID{&uid},
		Aliases: []*Alias{},
	}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Unexpected domain data. Got %#v, expected %#v", got, expected)
	}
}

func TestDeleteDomain(t *testing.T) {
	c := dmgr_setup()
	defer dmgr_teardown(c, t)

	mgr := DomainManager(&MongoDomainManager{c})
	uid := gouuid.New()
	mgr.ClaimDomain("gotest.org", &uid)
	mgr.ClaimDomain("gotest2.org", &uid)

	err := mgr.DeleteDomain("gotest.org", &uid)
	if err != nil {
		t.Fatalf("Could not delete domain: %s", err)
	}

	got, err := mgr.DomainsByOwner(&uid)
	if err != nil {
		t.Fatalf("Could not get remaining domains: %s", err)
	}
	expected := []*Domain{
		&Domain{
			Name:    "gotest2.org",
			Owners:  []*gouuid.UUID{&uid},
			Aliases: []*Alias{},
		},
	}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Unexpected domain data. Got %#v, expected %#v", got, expected)
	}
}

func TestDeleteUnknownDomain(t *testing.T) {
	c := dmgr_setup()
	defer dmgr_teardown(c, t)

	mgr := DomainManager(&MongoDomainManager{c})
	uid := gouuid.New()
	other_uid := gouuid.New()
	mgr.ClaimDomain("gotest.org", &uid)
	mgr.ClaimDomain("gotest2.org", &uid)

	err := mgr.DeleteDomain("gotest.org", &other_uid)
	if err == nil {
		t.Fatalf("Could delete domain: %s", err)
	}

	got, err := mgr.DomainsByOwner(&uid)
	if err != nil {
		t.Fatalf("Could not get remaining domains: %s", err)
	}
	if len(got) != 2 {
		t.Fatalf("Got unexpected number of remaining domains")
	}
}

func TestDomainsByOwner(t *testing.T) {
	c := dmgr_setup()
	defer dmgr_teardown(c, t)

	mgr := DomainManager(&MongoDomainManager{c})
	uid1, uid2 := gouuid.New(), gouuid.New()
	mgr.ClaimDomain("gotest_user1_1.org", &uid1)
	mgr.ClaimDomain("gotest_user1_2.org", &uid1)
	mgr.ClaimDomain("gotest_user2_1.org", &uid2)
	mgr.ClaimDomain("gotest_user2_2.org", &uid2)

	got, err := mgr.DomainsByOwner(&uid1)
	if err != nil {
		t.Fatalf("Could find domains: %s", err)
	}

	expected := []*Domain{
		&Domain{
			Name:    "gotest_user1_1.org",
			Owners:  []*gouuid.UUID{&uid1},
			Aliases: []*Alias{},
		},
		&Domain{
			Name:    "gotest_user1_2.org",
			Owners:  []*gouuid.UUID{&uid1},
			Aliases: []*Alias{},
		},
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Unexpected domain data. Got %#v, expected %#v", got, expected)
	}
}

func TestSetAlias(t *testing.T) {
	c := dmgr_setup()
	defer dmgr_teardown(c, t)

	mgr := DomainManager(&MongoDomainManager{c})
	uid := gouuid.New()
	mgr.ClaimDomain("gotest.org", &uid)

	expected := []*Alias{
		&Alias{
			RepoURL:    "repo1",
			RepoType:   "git",
			ForwardURL: "homepage1",
			Alias:      "alias1",
		},
		&Alias{
			RepoURL:    "repo2",
			RepoType:   "git",
			ForwardURL: "homepage2",
			Alias:      "alias2",
		},
	}
	mgr.SetAlias("gotest.org", expected[0], &uid)
	mgr.SetAlias("gotest.org", expected[1], &uid)

	domain := Domain{}
	c.Find(bson.M{}).One(&domain)
	if !reflect.DeepEqual(domain.Aliases, expected) {
		t.Fatalf("Unexpected domain data. Got %#v, expected %#v", domain.Aliases, expected)
	}
}

func TestWrongSetAlias(t *testing.T) {
	c := dmgr_setup()
	defer dmgr_teardown(c, t)

	mgr := DomainManager(&MongoDomainManager{c})
	uid := gouuid.New()
	other_uid := gouuid.New()
	mgr.ClaimDomain("gotest.org", &uid)

	expected := []*Alias{
		&Alias{
			RepoURL:    "repo1",
			RepoType:   "git",
			ForwardURL: "homepage1",
			Alias:      "alias1",
		},
		&Alias{
			RepoURL:    "repo2",
			RepoType:   "git",
			ForwardURL: "homepage2",
			Alias:      "alias2",
		},
	}
	mgr.SetAlias("gotest.org", expected[0], &uid)
	mgr.SetAlias("gotest.org", expected[1], &uid)
	err := mgr.SetAlias("gotest.org", &Alias{
		RepoURL:    "repo2",
		RepoType:   "git",
		ForwardURL: "homepage2",
		Alias:      "alias2",
	}, &other_uid)

	if err == nil {
		t.Fatalf("Could set alias with wrong id")
	}

	domain := Domain{}
	c.Find(bson.M{}).One(&domain)
	if !reflect.DeepEqual(domain.Aliases, expected) {
		t.Fatalf("Unexpected domain data. Got %#v, expected %#v", domain.Aliases, expected)
	}
}

func TestDeleteAlias(t *testing.T) {
	c := dmgr_setup()
	defer dmgr_teardown(c, t)

	mgr := DomainManager(&MongoDomainManager{c})
	uid := gouuid.New()
	mgr.ClaimDomain("gotest.org", &uid)

	expected := []*Alias{
		&Alias{
			RepoURL:    "repo1",
			RepoType:   "git",
			ForwardURL: "homepage1",
			Alias:      "alias1",
		},
		&Alias{
			RepoURL:    "repo2",
			RepoType:   "git",
			ForwardURL: "homepage2",
			Alias:      "alias2",
		},
	}
	mgr.SetAlias("gotest.org", expected[0], &uid)
	mgr.SetAlias("gotest.org", expected[1], &uid)

	mgr.DeleteAlias(expected[1].ID)
	expected = expected[0:1]

	domain := Domain{}
	c.Find(bson.M{}).One(&domain)
	if !reflect.DeepEqual(domain.Aliases, expected) {
		t.Fatalf("Unexpected domain data. Got %#v, expected %#v", domain.Aliases, expected)
	}
}

func TestFindAlias(t *testing.T) {
	c := dmgr_setup()
	defer dmgr_teardown(c, t)

	mgr := DomainManager(&MongoDomainManager{c})
	uid := gouuid.New()
	mgr.ClaimDomain("gotest.org", &uid)

	expected := []*Alias{
		&Alias{
			RepoURL:    "repo1",
			RepoType:   "git",
			ForwardURL: "homepage1",
			Alias:      "alias1",
		},
		&Alias{
			RepoURL:    "repo2",
			RepoType:   "git",
			ForwardURL: "homepage2",
			Alias:      "alias2",
		},
	}
	mgr.SetAlias("gotest.org", expected[0], &uid)
	mgr.SetAlias("gotest.org", expected[1], &uid)

	alias, err := mgr.FindAlias("gotest.org", "alias1")
	if err != nil {
		t.Fatalf("Could not find domain: %s", err)
	}

	if !reflect.DeepEqual(alias, expected[0]) {
		t.Fatalf("Unexpected domain data. Got %#v, expected %#v", alias, expected)
	}
}
