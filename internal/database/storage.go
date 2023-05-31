package database

import (
	"sync"
)

type MemStorage struct {
	// login-pass pairs stored for each user.
	Logins map[string][]LoginCreds

	// notes stored for each user
	Notes map[string][]Note

	//cards stored for each user
	Cards map[string][]Card

	//Binary data stored for each user
	Binaries map[string][]Binary

	// to ensure possible concurrent usage
	Mu *sync.Mutex
}

func NewStorage() *MemStorage {
	return &MemStorage{
		Logins:   make(map[string][]LoginCreds),
		Notes:    make(map[string][]Note),
		Cards:    make(map[string][]Card),
		Binaries: make(map[string][]Binary),
		Mu:       &sync.Mutex{},
	}
}

func (s *MemStorage) ListLogins(username string) ([]LoginCreds, error) {
	return s.Logins[username], nil
}

func (s *MemStorage) ListCards(username string) ([]Card, error) {
	return s.Cards[username], nil
}

func (s *MemStorage) ListNotes(username string) ([]Note, error) {
	return s.Notes[username], nil
}

func (s *MemStorage) ListBinaries(username string) ([]Binary, error) {
	return s.Binaries[username], nil
}

func (s *MemStorage) SetBinary(username string, newbinary Binary) error {
	if _, ok := s.Binaries[username]; !ok {
		s.Binaries[username] = make([]Binary, 1)
	}
	s.Binaries[username] = append(s.Binaries[username], newbinary)
	return nil
}

func (s *MemStorage) GetBinary(username string, binaryname string) (Binary, error) {

	for _, bindata := range s.Binaries[username] {
		if bindata.Name == binaryname {
			return bindata, nil
		}
	}

	return Binary{}, ErrDataNotFound
}

func (s *MemStorage) SetCard(username string, card Card) error {
	if _, ok := s.Cards[username]; !ok {
		s.Cards[username] = make([]Card, 1)
	}
	s.Cards[username] = append(s.Cards[username], card)
	return nil
}

func (s *MemStorage) SetLoginCred(username string, logindata LoginCreds) error {
	if _, ok := s.Logins[username]; !ok {
		s.Logins[username] = make([]LoginCreds, 1)
	}
	s.Logins[username] = append(s.Logins[username], logindata)
	return nil
}

func (s *MemStorage) GetLoginCred(username string, loginname string) (LoginCreds, error) {
	for _, logindata := range s.Logins[username] {
		if logindata.Name == loginname {
			return logindata, nil
		}
	}
	return LoginCreds{}, ErrDataNotFound
}

func (s *MemStorage) GetCard(username string, cardname string) (Card, error) {
	for _, carddata := range s.Cards[username] {
		if carddata.Name == cardname {
			return carddata, nil
		}
	}
	return Card{}, ErrDataNotFound
}

func (s *MemStorage) ListLoginCreds(username string) ([]LoginCreds, error) {
	if _, ok := s.Logins[username]; !ok {
		return nil, ErrDataNotFound
	}
	return s.Logins[username], nil
}
