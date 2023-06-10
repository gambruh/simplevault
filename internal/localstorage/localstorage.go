package localstorage

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/gambruh/gophkeeper/internal/config"
	"github.com/gambruh/gophkeeper/internal/database"
	"github.com/gambruh/gophkeeper/internal/encrypt"
)

type LocalStorage struct {
	Cards      []string
	Logincreds []string
	Notes      []string
	Binaries   []string
	Mu         sync.Mutex
}

const (
	cardsFile      = "/cards"
	notesFile      = "/notes"
	binariesFile   = "/binaries"
	loginCredsFile = "/logincred"
)

var (
	ErrMetanameIsTaken = errors.New("metaname is taken, please provide another item name")
	ErrNoData          = errors.New("data not found")
)

func NewStorage() *LocalStorage {
	ls := &LocalStorage{
		Cards:      make([]string, 0),
		Logincreds: make([]string, 0),
		Notes:      make([]string, 0),
		Binaries:   make([]string, 0),
		Mu:         sync.Mutex{},
	}

	return ls

}

func (s *LocalStorage) InitStorage(key []byte) error {
	list, err := s.ListCardsFromFile(key)
	if err != nil {
		s.deleteCardsFile()
		s.Cards = []string{}
	} else {
		s.Cards = list
	}

	list, err = s.ListLoginCredsFromFile(key)
	if err != nil {
		s.deleteLoginCredsFile()
		s.Logincreds = []string{}
		return nil
	} else {
		s.Logincreds = list
	}

	return nil
}

func (s *LocalStorage) DeleteLocalStorage() error {

	if err := s.deleteCardsFile(); err != nil {
		return err
	}

	if err := s.deleteLoginCredsFile(); err != nil {
		return err
	}

	return nil
}

func (s *LocalStorage) deleteCardsFile() error {
	err := os.Remove(config.ClientCfg.LocalStorage + cardsFile)
	if err != nil {
		return fmt.Errorf("can't delete local cache:%w", err)
	}
	return nil
}

func (s *LocalStorage) deleteLoginCredsFile() error {
	err := os.Remove(config.ClientCfg.LocalStorage + loginCredsFile)
	if err != nil {
		return fmt.Errorf("can't delete local cache:%w", err)
	}
	return nil
}

// SaveCard method encrypts and saves card data to the storage
func (s *LocalStorage) SaveCard(card database.Card, key []byte) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	// check if the card with this name is in storage. Return error if yes
	if check := s.lookupCard(card.Cardname); check {
		return ErrMetanameIsTaken
	}

	// add cardname to check array
	s.Cards = append(s.Cards, card.Cardname)

	file, err := os.OpenFile(config.ClientCfg.LocalStorage+cardsFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("error in SaveCard when opening file:%w", err)
	}
	defer file.Close()

	// concatenating card to string
	cardStr := card.Cardname + "," + card.Number + "," + card.Name + "," + card.Surname + "," + card.ValidTill + "," + card.Code

	// encrypting the card data
	encrypted, err := encrypt.EncryptData([]byte(cardStr), key)
	if err != nil {
		return err
	}
	// Encode the encrypted password in base64 for storage
	encodedData := base64.StdEncoding.EncodeToString(encrypted)

	_, err = fmt.Fprintf(file, "%s\n", encodedData)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return fmt.Errorf("error in SaveCard when writing in file:%w", err)
	}

	return nil
}

func (s *LocalStorage) GetCard(cardname string, key []byte) (card database.Card, err error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	// check if the card with this name is in storage. Return error if yes
	if check := s.lookupCard(cardname); !check {
		return database.Card{}, ErrNoData
	}

	// opening the localstorage file
	file, err := os.OpenFile(config.ClientCfg.LocalStorage+cardsFile, os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		return database.Card{}, fmt.Errorf("error in GetCard when opening file:%w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// reading with Scanner each line, until encounter needed one
	for scanner.Scan() {
		line := scanner.Text()
		dst, err := base64.StdEncoding.DecodeString(line)
		if err != nil {
			return database.Card{}, err
		}

		decryptedData, err := encrypt.DecryptData(dst, key)
		if err != nil {
			return database.Card{}, err
		}

		//Getting the card string, splitting it by comma to get values
		cardStr := string(decryptedData)
		cardArr := strings.Split(cardStr, ",")

		// cardArr[0] is the cardname. If it is the one we are looking for,
		// we fill in fields of database.Card struct and return it
		if cardArr[0] == cardname {
			card.Cardname = cardArr[0]
			card.Number = cardArr[1]
			card.Name = cardArr[2]
			card.Surname = cardArr[3]
			card.ValidTill = cardArr[4]
			card.Code = cardArr[5]
			return card, nil
		}
	}
	return database.Card{}, database.ErrDataNotFound
}

func (s *LocalStorage) ListCards() (cards []string, err error) {
	return s.Cards, nil
}

func (s *LocalStorage) ListCardsFromFile(key []byte) (cards []string, err error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	// opening the localstorage file
	file, err := os.OpenFile(config.ClientCfg.LocalStorage+cardsFile, os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		return nil, fmt.Errorf("error in ListCards when opening file:%w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		dst, err := base64.StdEncoding.DecodeString(line)
		if err != nil {
			return nil, err
		}
		decryptedData, err := encrypt.DecryptData(dst, key)
		if err != nil {
			return nil, err
		}

		//Getting the card string, splitting it by comma to get values
		cardStr := string(decryptedData)
		cardArr := strings.Split(cardStr, ",")

		// cardArr[0] is the cardname
		cards = append(cards, cardArr[0])
	}

	if len(cards) == 0 {
		return nil, ErrNoData
	}

	return cards, nil
}

func (s *LocalStorage) lookupCard(cardname string) bool {

	for _, c := range s.Cards {
		if c == cardname {
			return true
		}
	}
	return false
}

func (s *LocalStorage) SaveLoginCreds(logincreds database.LoginCreds, key []byte) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	// check if the card with this name is in storage. Return error if yes
	if check := s.lookupLoginCreds(logincreds.Name); check {
		return ErrMetanameIsTaken
	}

	// add name to check array
	s.Logincreds = append(s.Logincreds, logincreds.Name)

	file, err := os.OpenFile(config.ClientCfg.LocalStorage+loginCredsFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("error in SaveCard when opening file:%w", err)
	}
	defer file.Close()

	// concatenating string
	logincredsStr := logincreds.Name + "," + logincreds.Site + "," + logincreds.Login + "," + logincreds.Password

	// encrypting the data
	encrypted, err := encrypt.EncryptData([]byte(logincredsStr), key)
	if err != nil {
		return err
	}
	// encoding the encrypted data in base64 for storage
	encodedData := base64.StdEncoding.EncodeToString(encrypted)

	// saving data to the filestorage
	_, err = fmt.Fprintf(file, "%s\n", encodedData)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return fmt.Errorf("error in SaveLoginCreds when writing in file:%w", err)
	}

	return nil
}

func (s *LocalStorage) lookupLoginCreds(logincreds string) bool {

	for _, l := range s.Logincreds {
		if l == logincreds {
			return true
		}
	}
	return false
}

func (s *LocalStorage) GetLoginCreds(logincredsname string, key []byte) (logincreds database.LoginCreds, err error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	// check if the card with this name is in storage. Return error if yes
	if check := s.lookupLoginCreds(logincredsname); !check {
		return database.LoginCreds{}, ErrNoData
	}

	// opening the localstorage file
	file, err := os.OpenFile(config.ClientCfg.LocalStorage+loginCredsFile, os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		return database.LoginCreds{}, fmt.Errorf("error in GetLoginCreds when opening file:%w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// reading with Scanner each line, until encounter needed one
	for scanner.Scan() {
		line := scanner.Text()
		dst, err := base64.StdEncoding.DecodeString(line)
		if err != nil {
			return database.LoginCreds{}, err
		}

		decryptedData, err := encrypt.DecryptData(dst, key)
		if err != nil {
			return database.LoginCreds{}, err
		}

		//Getting the string, splitting it by comma to get values
		loginCredStr := string(decryptedData)
		loginCredArr := strings.Split(loginCredStr, ",")

		if loginCredArr[0] == logincredsname {
			logincreds.Name = loginCredArr[0]
			logincreds.Site = loginCredArr[1]
			logincreds.Login = loginCredArr[2]
			logincreds.Password = loginCredArr[3]

			return logincreds, nil
		}
	}
	return database.LoginCreds{}, database.ErrDataNotFound
}

func (s *LocalStorage) ListLoginCreds() (logincreds []string, err error) {
	return s.Logincreds, nil
}

func (s *LocalStorage) ListLoginCredsFromFile(key []byte) (logincreds []string, err error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	// opening the localstorage file
	file, err := os.OpenFile(config.ClientCfg.LocalStorage+loginCredsFile, os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		return nil, fmt.Errorf("error in ListLoginCredsFromFile when opening file:%w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		dst, err := base64.StdEncoding.DecodeString(line)
		if err != nil {
			return nil, err
		}
		decryptedData, err := encrypt.DecryptData(dst, key)
		if err != nil {
			return nil, err
		}

		//Getting the card string, splitting it by comma to get values
		loginCredStr := string(decryptedData)
		loginCredArr := strings.Split(loginCredStr, ",")

		// cardArr[0] is the cardname
		logincreds = append(logincreds, loginCredArr[0])
	}

	if len(logincreds) == 0 {
		return nil, ErrNoData
	}

	return logincreds, nil
}

// Notes processing methods
func (s *LocalStorage) SaveNote(note database.Note, key []byte) error {

	return nil
}

func (s *LocalStorage) GetNote(notename string, key []byte) (note database.Note, err error) {

	return database.Note{}, nil
}

func (s *LocalStorage) ListNotes() (notes []string, err error) {
	return s.Notes, nil
}

// Binaries processing methods

func (s *LocalStorage) SaveBinary(binary database.Binary, key []byte) error {
	return nil
}

func (s *LocalStorage) GetBinary(binaryname string, key []byte) (binary database.Binary, err error) {
	return database.Binary{}, nil
}

func (s *LocalStorage) ListBinaries() (binaries []string, err error) {
	return s.Binaries, nil
}
