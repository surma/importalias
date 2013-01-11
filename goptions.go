package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gorilla/sessions"
)

type SessionStore struct {
	sessions.Store
}

func (s *SessionStore) MarshalGoption(key string) error {
	if len([]byte(key)) != 32 {
		return fmt.Errorf("Cookie key needs to be 32 byte")
	}
	cs := sessions.NewCookieStore([]byte(key))
	cs.Options.MaxAge = int(options.SessionTTL / time.Second)
	s.Store = cs
	return nil
}

type AuthConfig struct {
	Type      string `json:"type"`
	AuthURL   string `json:"auth_url"`
	TokenURL  string `json:"token_url"`
	Scope     string `json:"scope"`
	AuthKey   *AuthKey
	Extractor struct {
		Type  string `json:"type"`
		URL   string `json:"url"`
		Field string `json:"field"`
	} `json:"extractor"`
}

type AuthList map[string]*AuthConfig

func (a *AuthList) MarshalGoption(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("Could not open file: %s", err)
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(a)
	if err != nil {
		return fmt.Errorf("Could not decode file: %s", err)
	}
	return nil
}

type AuthKey struct {
	Name     string
	ClientID string
	Secret   string
}

func (a *AuthKey) MarshalGoption(key string) error {
	keyparts := strings.Split(key, ":")
	if len(keyparts) < 3 {
		return fmt.Errorf("Invalid auth key format \"%s\"", key)
	}
	name, clientid, secret := keyparts[0], keyparts[1], keyparts[2]
	a.Name = name
	a.ClientID = clientid
	a.Secret = secret
	return nil
}
