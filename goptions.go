package main

import (
	"encoding/json"
	"fmt"
	"os"
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
