package database

import (
	"fmt"
	"sync"
)

type MemStorage struct {
	// login-pass pairs stored for each user.
	Logins map[string][]string

	// notes stored for each user
	Notes map[string][]string

	//cards stored for each user
	Cards map[string][]string

	//Binary data stored for each user
	Binaries map[string][]string

	// to ensure possible concurrent usage
	Mu *sync.Mutex
}

func NewStorage() *MemStorage {
	return &MemStorage{
		Logins:   make(map[string][]string),
		Notes:    make(map[string][]string),
		Cards:    make(map[string][]string),
		Binaries: make(map[string][]string),
		Mu:       &sync.Mutex{},
	}
}

func (s *MemStorage) ListLogins(username string) ([]string, error) {
	return s.Logins[username], nil
}

func (s *MemStorage) ListCards(username string) ([]string, error) {
	return s.Cards[username], nil
}

func (s *MemStorage) ListNotes(username string) ([]string, error) {
	return s.Notes[username], nil
}

func (s *MemStorage) ListBinaries(username string) ([]string, error) {
	return s.Binaries[username], nil
}

func (s *MemStorage) SetBinary(username string, newbinary string) error {
	if _, ok := s.Binaries[username]; !ok {
		s.Binaries[username] = make([]string, 1)
	}
	s.Binaries[username] = append(s.Binaries[username], newbinary)
	return nil
}

func (s *MemStorage) GetBinary(username string, binaryname string) (Binary, error) {

	for _, bindata := range s.Binaries[username] {
		fmt.Println(bindata)
	}

	return Binary{}, ErrDataNotFound
}

func (s *MemStorage) SetCard(username string, card string) error {
	if _, ok := s.Cards[username]; !ok {
		s.Cards[username] = make([]string, 1)
	}
	s.Cards[username] = append(s.Cards[username], card)
	return nil
}

func (s *MemStorage) SetLoginCred(username string, logindata string) error {
	if _, ok := s.Logins[username]; !ok {
		s.Logins[username] = make([]string, 1)
	}
	s.Logins[username] = append(s.Logins[username], logindata)
	return nil
}

func (s *MemStorage) GetLoginCred(username string, loginname string) (LoginCreds, error) {
	for _, logindata := range s.Logins[username] {

		fmt.Println(logindata)

	}
	return LoginCreds{}, ErrDataNotFound
}

func (s *MemStorage) GetCard(username string, cardname string) (Card, error) {
	for _, carddata := range s.Cards[username] {

		fmt.Print(carddata)

	}
	return Card{}, ErrDataNotFound
}

func (s *MemStorage) ListLoginCreds(username string) ([]string, error) {
	if _, ok := s.Logins[username]; !ok {
		return nil, ErrDataNotFound
	}
	return s.Logins[username], nil
}
