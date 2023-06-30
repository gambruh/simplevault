package storage

import "errors"

type Storage interface {
	SetLoginCred(username string, logincreds EncryptedData) error
	SetNote(username string, note EncryptedData) error
	SetBinary(username string, binary Binary) error
	SetCard(username string, card EncryptedData) error
	GetLoginCred(username string, name string) (EncryptedData, error)
	GetNote(username string, name string) (EncryptedData, error)
	GetBinary(username string, name string) (Binary, error)
	GetCard(username string, name string) (EncryptedData, error)
	ListLoginCreds(username string) ([]string, error)
	ListNotes(username string) ([]string, error)
	ListBinaries(username string) ([]string, error)
	ListCards(username string) ([]string, error)
}

type LoginCreds struct {
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Site     string `json:"site"`
}

type Note struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

type Binary struct {
	Name string `json:"name"`
	Data []byte `json:"data"`
}

type Card struct {
	Cardname  string `json:"cardname"`
	Number    string `json:"number"`
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	ValidTill string `json:"valid till"`
	Code      string `json:"code"`
}

type EncryptedData struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// errors
var (
	ErrTableDoesntExist = errors.New("table doesn't exist")
	ErrWrongPassword    = errors.New("wrong password")
	ErrDataNotFound     = errors.New("requested data not found in storage")
	ErrMetanameIsTaken  = errors.New("metaname is already in use")
)
