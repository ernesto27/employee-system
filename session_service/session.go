package session_service

import (
	"net/http"

	"github.com/gorilla/sessions"
)

type Session struct {
	store        *sessions.CookieStore
	loginName    string
	autenticated string
}

func NewSession() *Session {
	var store = sessions.NewCookieStore([]byte("your-secret-key"))
	return &Session{
		store:        store,
		loginName:    "loggin",
		autenticated: "authenticated",
	}
}

func (s *Session) CreateLogin(w http.ResponseWriter, r *http.Request) error {
	session, _ := s.store.Get(r, s.loginName)
	session.Values[s.autenticated] = true
	err := session.Save(r, w)

	if err != nil {
		return err
	}

	return nil
}

func (s *Session) RemoveLogin(w http.ResponseWriter, r *http.Request) error {
	session, _ := s.store.Get(r, s.loginName)
	session.Values[s.autenticated] = false
	err := session.Save(r, w)

	if err != nil {
		return err
	}

	return nil
}

func (s *Session) CheckLogin(r *http.Request) bool {
	session, _ := s.store.Get(r, s.loginName)
	if session.Values[s.autenticated] == nil || session.Values[s.autenticated] == false {
		return false
	}

	return true
}
