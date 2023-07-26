// Package memstorage provides inmemory implementation of Storage interface
// Purpose of this package is to provide struct and methods to use in unit-tests mainly
package memstorage

import (
	"fmt"
	"sync"

	"github.com/gambruh/simplevault/internal/storage"
)

// MemStorage struct is supposed to be used in unit-tests only
type MemStorage struct {
	// login-pass pairs stored for each user.
	Logins map[string][]storage.LoginCreds

	// notes stored for each user
	Notes map[string][]storage.Note

	//cards stored for each user
	Cards map[string][]storage.Card

	//Binary data stored for each user
	Binaries map[string][]storage.Binary

	// to ensure possible concurrent usage
	Mu *sync.Mutex
}

// NewStorage is a constructor of a new MemStorage struct
func NewStorage() *MemStorage {
	return &MemStorage{
		Logins:   make(map[string][]storage.LoginCreds),
		Notes:    make(map[string][]storage.Note),
		Cards:    make(map[string][]storage.Card),
		Binaries: make(map[string][]storage.Binary),
		Mu:       &sync.Mutex{},
	}
}

// ListLoginCreds returns a list of names of login credentials saved in Storage
func (s *MemStorage) ListLoginCreds(username string) ([]string, error) {
	var list []string

	for _, item := range s.Logins[username] {
		list = append(list, item.Name)
	}

	return list, nil
}

// ListCards returns a list of names of cards saved in Storage
func (s *MemStorage) ListCards(username string) ([]string, error) {
	var list []string

	for _, item := range s.Cards[username] {
		list = append(list, item.Name)
	}

	return list, nil
}

// ListNotes returns a list of names of notes saved in Storage
func (s *MemStorage) ListNotes(username string) ([]string, error) {
	var list []string

	for _, item := range s.Notes[username] {
		list = append(list, item.Name)
	}

	return list, nil
}

// ListBinaries returns a list of names of binaries saved in Storage
func (s *MemStorage) ListBinaries(username string) ([]string, error) {
	var list []string

	for _, item := range s.Binaries[username] {
		list = append(list, item.Name)
	}

	return list, nil
}

// SetBinary saves a new binary in the Storage
func (s *MemStorage) SetBinary(username string, newbinary storage.Binary) error {
	if _, ok := s.Binaries[username]; !ok {
		s.Binaries[username] = make([]storage.Binary, 1)
	}
	s.Binaries[username] = append(s.Binaries[username], newbinary)
	return nil
}

// GetBinary returns a binary by it's name
func (s *MemStorage) GetBinary(username string, binaryname string) (storage.Binary, error) {

	for _, bindata := range s.Binaries[username] {
		fmt.Println(bindata)
	}

	return storage.Binary{}, storage.ErrDataNotFound
}

// SetCard saves a card in the Storage
func (s *MemStorage) SetCard(username string, card storage.Card) error {
	if _, ok := s.Cards[username]; !ok {
		s.Cards[username] = make([]storage.Card, 1)
	}
	s.Cards[username] = append(s.Cards[username], card)
	return nil
}

// SetLoginCred saves a login credentials in the Storage
func (s *MemStorage) SetLoginCred(username string, logindata storage.LoginCreds) error {
	if _, ok := s.Logins[username]; !ok {
		s.Logins[username] = make([]storage.LoginCreds, 1)
	}
	s.Logins[username] = append(s.Logins[username], logindata)
	return nil
}

// GetLoginCred returns a login credentials by it's name
func (s *MemStorage) GetLoginCred(username string, loginname string) (storage.LoginCreds, error) {
	for _, logindata := range s.Logins[username] {

		fmt.Println(logindata)

	}
	return storage.LoginCreds{}, storage.ErrDataNotFound
}

// GetCard returns a card by it's name
func (s *MemStorage) GetCard(username string, cardname string) (storage.Card, error) {
	for _, carddata := range s.Cards[username] {

		fmt.Print(carddata)

	}
	return storage.Card{}, storage.ErrDataNotFound
}
