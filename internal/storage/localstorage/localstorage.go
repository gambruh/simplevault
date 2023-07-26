// Package localstorage provides filestorage implementation of Storage interface for a local client's storage
package localstorage

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gambruh/gophkeeper/internal/config"
	"github.com/gambruh/gophkeeper/internal/encrypt"
	"github.com/gambruh/gophkeeper/internal/storage"
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
	loginCredsFile = "/logincred"
	binariesFolder = "/binaries"
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

// InitStorage creates required directories if they doesn't exist yet
// If files and directories exists, it checks for it's contents and get the data loaded into struct fields
func (s *LocalStorage) InitStorage(key []byte) error {
	//create local folders if needed
	os.Mkdir(config.ClientCfg.LocalStorage, 0600)
	os.Mkdir(config.ClientCfg.BinInputFolder, 0600)
	os.Mkdir(config.ClientCfg.BinOutputFolder, 0600)

	list, err := s.ListCardsFromFile(key)
	if err != nil {
		s.Cards = []string{}
	} else {
		s.Cards = list
	}

	list, err = s.ListLoginCredsFromFile(key)
	if err != nil {
		s.Logincreds = []string{}
	} else {
		s.Logincreds = list
	}

	list, err = s.ListNotesFromFile(key)
	if err != nil {
		s.Notes = []string{}
	} else {
		s.Notes = list
	}
	list, err = s.ListBinariesFromFolder()
	if err != nil {
		s.Binaries = []string{}
	} else {
		s.Binaries = list
	}

	return nil

}

// DeleteLocalStorage removes files from local file storage
func (s *LocalStorage) DeleteLocalStorage() error {

	if err := s.deleteCardsFile(); err != nil {
		return err
	}

	if err := s.deleteLoginCredsFile(); err != nil {
		return err
	}

	if err := s.deleteNotesFile(); err != nil {
		return err
	}

	if err := s.deleteBinaryFiles(); err != nil {
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

func (s *LocalStorage) deleteNotesFile() error {
	err := os.Remove(config.ClientCfg.LocalStorage + notesFile)
	if err != nil {
		return fmt.Errorf("can't delete local cache:%w", err)
	}
	return nil
}

func (s *LocalStorage) deleteBinaryFiles() error {
	err := os.RemoveAll(config.ClientCfg.LocalStorage + binariesFolder)
	if err != nil {
		return fmt.Errorf("can't delete local cache:%w", err)
	}
	return nil
}

// SaveCard method encrypts and saves card data to the storage
func (s *LocalStorage) SaveCard(card storage.Card, key []byte) error {
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

// GetCard gets the card data by its name, decrypts it and returns as a Card struct
func (s *LocalStorage) GetCard(cardname string, key []byte) (card storage.Card, err error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	// check if the card with this name is in storage. Return error if yes
	if check := s.lookupCard(cardname); !check {
		return storage.Card{}, ErrNoData
	}

	// opening the localstorage file
	file, err := os.OpenFile(config.ClientCfg.LocalStorage+cardsFile, os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		return storage.Card{}, fmt.Errorf("error in GetCard when opening file:%w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// reading with Scanner each line, until encounter needed one
	for scanner.Scan() {
		line := scanner.Text()
		dst, err := base64.StdEncoding.DecodeString(line)
		if err != nil {
			return storage.Card{}, err
		}

		decryptedData, err := encrypt.DecryptData(dst, key)
		if err != nil {
			return storage.Card{}, err
		}

		//Getting the card string, splitting it by comma to get values
		cardStr := string(decryptedData)
		cardArr := strings.Split(cardStr, ",")

		// cardArr[0] is the cardname. If it is the one we are looking for,
		// we fill in fields of storage.Card struct and return it
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
	return storage.Card{}, storage.ErrDataNotFound
}

// ListCards returns a list of names of cards saved in Storage.
// Names are taken from s.Cards field of the LocalStorage struct
func (s *LocalStorage) ListCards() (cards []string, err error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	if len(s.Cards) == 0 {
		return nil, nil
	}
	return s.Cards, nil
}

// ListCardsFromFile checks cards file and returns a list of names of cards saved in it.
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

// SaveLoginCreds encrypts and saves login credentials in the client's file storage
func (s *LocalStorage) SaveLoginCreds(logincreds storage.LoginCreds, key []byte) error {
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
		return fmt.Errorf("error in SaveLoginCreds when opening file:%w", err)
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

// GetLoginCreds gets login credentials from local storage by the logincreds name
func (s *LocalStorage) GetLoginCreds(logincredsname string, key []byte) (logincreds storage.LoginCreds, err error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	// check if the card with this name is in storage. Return error if yes
	if check := s.lookupLoginCreds(logincredsname); !check {
		return storage.LoginCreds{}, ErrNoData
	}

	// opening the localstorage file
	file, err := os.OpenFile(config.ClientCfg.LocalStorage+loginCredsFile, os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		return storage.LoginCreds{}, fmt.Errorf("error in GetLoginCreds when opening file:%w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// reading with Scanner each line, until encounter needed one
	for scanner.Scan() {
		line := scanner.Text()
		dst, err := base64.StdEncoding.DecodeString(line)
		if err != nil {
			return storage.LoginCreds{}, err
		}

		decryptedData, err := encrypt.DecryptData(dst, key)
		if err != nil {
			return storage.LoginCreds{}, err
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
	return storage.LoginCreds{}, storage.ErrDataNotFound
}

// ListLoginCreds returns a list of login credential names saves in a structure
func (s *LocalStorage) ListLoginCreds() (logincreds []string, err error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	if len(s.Logincreds) == 0 {
		return nil, nil
	}
	return s.Logincreds, nil
}

// ListLoginCredsFromFile checks and returns list of login credential names saves from the file
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

func (s *LocalStorage) lookupNote(note string) bool {

	for _, l := range s.Notes {
		if l == note {
			return true
		}
	}
	return false
}

// Notes processing methods
func (s *LocalStorage) SaveNote(note storage.Note, key []byte) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	// check if the card with this name is in storage. Return error if yes
	if check := s.lookupNote(note.Name); check {
		return ErrMetanameIsTaken
	}

	// add name to check array
	s.Notes = append(s.Notes, note.Name)

	file, err := os.OpenFile(config.ClientCfg.LocalStorage+notesFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("error in SaveNote when opening file:%w", err)
	}
	defer file.Close()

	// concatenating string
	noteStr := note.Name + "," + note.Text

	// encrypting the data
	encrypted, err := encrypt.EncryptData([]byte(noteStr), key)
	if err != nil {
		return err
	}
	// encoding the encrypted data in base64 for storage
	encodedData := base64.StdEncoding.EncodeToString(encrypted)

	// saving data to the filestorage
	_, err = fmt.Fprintf(file, "%s\n", encodedData)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return fmt.Errorf("error in SaveNote when writing in file:%w", err)
	}

	return nil
}

func (s *LocalStorage) GetNote(notename string, key []byte) (note storage.Note, err error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	// check if the card with this name is in storage. Return error if yes
	if check := s.lookupNote(notename); !check {
		return storage.Note{}, ErrNoData
	}

	// opening the localstorage file
	file, err := os.OpenFile(config.ClientCfg.LocalStorage+notesFile, os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		return storage.Note{}, fmt.Errorf("error in GetNote when opening file:%w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// reading with Scanner each line, until encounter needed one
	for scanner.Scan() {
		line := scanner.Text()
		dst, err := base64.StdEncoding.DecodeString(line)
		if err != nil {
			return storage.Note{}, err
		}

		decryptedData, err := encrypt.DecryptData(dst, key)
		if err != nil {
			return storage.Note{}, err
		}

		//Getting the string, splitting it by comma to get values
		noteStr := string(decryptedData)
		noteArr := strings.Split(noteStr, ",")

		if noteArr[0] == notename {
			note.Name = noteArr[0]
			note.Text = noteArr[1]
			return note, nil
		}
	}
	return storage.Note{}, storage.ErrDataNotFound
}

func (s *LocalStorage) ListNotes() (notes []string, err error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	if len(s.Notes) == 0 {
		return nil, nil
	}
	return s.Notes, nil
}

func (s *LocalStorage) ListNotesFromFile(key []byte) (notes []string, err error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	// opening the localstorage file
	file, err := os.OpenFile(config.ClientCfg.LocalStorage+notesFile, os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		return nil, fmt.Errorf("error in ListNotesFromFile when opening file:%w", err)
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
		noteStr := string(decryptedData)
		noteArr := strings.Split(noteStr, ",")

		// cardArr[0] is the cardname
		notes = append(notes, noteArr[0])
	}

	if len(notes) == 0 {
		return nil, ErrNoData
	}

	return notes, nil
}

// Binaries processing methods
func (s *LocalStorage) lookupBinary(binaryname string) bool {
	for _, b := range s.Binaries {
		if b == binaryname {
			return true
		}
	}
	return false
}

func (s *LocalStorage) SaveBinary(binary storage.Binary, key []byte) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	// check if the card with this name is in storage. Return error if yes
	if check := s.lookupBinary(binary.Name); check {
		return ErrMetanameIsTaken
	}

	// add name to check array
	s.Binaries = append(s.Binaries, binary.Name)

	// just in case create binaries folder
	os.Mkdir(config.ClientCfg.LocalStorage+binariesFolder, 0600)

	file, err := os.OpenFile(config.ClientCfg.LocalStorage+binariesFolder+"/"+binary.Name, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("error in SaveBinary when opening file:%w", err)
	}
	defer file.Close()

	// encrypting the data
	encrypted, err := encrypt.EncryptData(binary.Data, key)
	if err != nil {
		return err
	}

	encodedData := base64.StdEncoding.EncodeToString(encrypted)

	// saving data to the filestorage
	_, err = fmt.Fprint(file, encodedData)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return fmt.Errorf("error in SaveBinary when writing in file:%w", err)
	}

	return nil
}

func (s *LocalStorage) GetBinary(binaryname string, key []byte) (binary storage.Binary, err error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	// check if the card with this name is in storage. Return error if yes
	if check := s.lookupBinary(binaryname); !check {
		return storage.Binary{}, ErrNoData
	}

	file, err := os.OpenFile(config.ClientCfg.LocalStorage+binariesFolder+"/"+binaryname, os.O_RDONLY, 0600)
	if err != nil {
		return storage.Binary{}, fmt.Errorf("error in GetBinary when opening file:%w", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	data, err := io.ReadAll(reader)

	if err != nil {
		return storage.Binary{}, err
	}

	dst, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return storage.Binary{}, err
	}
	decryptedData, err := encrypt.DecryptData(dst, key)
	if err != nil {
		return storage.Binary{}, err
	}

	binary.Name = binaryname
	binary.Data = decryptedData

	//creates the directory if its not there
	os.Mkdir(config.ClientCfg.BinOutputFolder, 0600)

	newfile, err := os.OpenFile(config.ClientCfg.BinOutputFolder+"/"+binary.Name, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return storage.Binary{}, fmt.Errorf("error in GetBinary when opening file:%w", err)
	}
	defer newfile.Close()

	// saving data to the filestorage
	_, err = fmt.Fprint(newfile, binary.Data)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return storage.Binary{}, fmt.Errorf("error in SaveBinary when writing in file:%w", err)
	}

	return binary, nil
}

func (s *LocalStorage) GetEncryptedBinary(binaryname string) (binary storage.Binary, err error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	// check if the card with this name is in storage. Return error if yes
	if check := s.lookupBinary(binaryname); !check {
		return storage.Binary{}, ErrNoData
	}

	file, err := os.OpenFile(config.ClientCfg.LocalStorage+binariesFolder+"/"+binaryname, os.O_RDONLY, 0600)
	if err != nil {
		return storage.Binary{}, fmt.Errorf("error in GetEncryptedBinary when opening file:%w", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	data, err := io.ReadAll(reader)
	if err != nil {
		return storage.Binary{}, err
	}

	binary.Name = binaryname
	binary.Data = data

	return binary, nil
}

func (s *LocalStorage) ListBinaries() (binaries []string, err error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	if len(s.Binaries) == 0 {
		return nil, nil
	}
	return s.Binaries, nil
}

func (s *LocalStorage) ListBinariesFromFolder() (binaries []string, err error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	// Make a folder if its a first time
	os.Mkdir(config.ClientCfg.LocalStorage+binariesFolder, 0600)

	// Open the folder
	err = filepath.Walk(config.ClientCfg.LocalStorage+binariesFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error accessing path:", err)
			return nil
		}

		if !info.IsDir() {
			binaries = append(binaries, info.Name())
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error walking folder:", err)
		return nil, err
	}

	return binaries, nil
}
