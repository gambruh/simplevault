package storage

import (
	"fmt"
	"sync"
)

// MemStorage struct is supposed to be used in unit-tests only
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

func (s *MemStorage) ListLoginCreds(username string) ([]string, error) {
	var list []string

	for _, item := range s.Logins[username] {
		list = append(list, item.Name)
	}

	return list, nil
}

func (s *MemStorage) ListCards(username string) ([]string, error) {
	var list []string

	for _, item := range s.Cards[username] {
		list = append(list, item.Name)
	}

	return list, nil
}

func (s *MemStorage) ListNotes(username string) ([]string, error) {
	var list []string

	for _, item := range s.Notes[username] {
		list = append(list, item.Name)
	}

	return list, nil
}

func (s *MemStorage) ListBinaries(username string) ([]string, error) {
	var list []string

	for _, item := range s.Binaries[username] {
		list = append(list, item.Name)
	}

	return list, nil
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
		fmt.Println(bindata)
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
