package localstorage

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/gambruh/gophkeeper/internal/config"
	"github.com/gambruh/gophkeeper/internal/database"
	"github.com/gambruh/gophkeeper/internal/encrypt"
)

type LocalStorage struct {
	Cards      string
	Logincreds string
	Notes      string
	Binaries   string
}

const (
	cardsFile      = "/cards"
	notesFile      = "/notes"
	binariesFile   = "/binaries"
	loginCredsFile = "/logincred"
)

func NewStorage() *LocalStorage {
	return &LocalStorage{
		Cards:      cardsFile,
		Logincreds: loginCredsFile,
		Notes:      notesFile,
		Binaries:   binariesFile,
	}
}

func (s *LocalStorage) SaveCard(card database.Card, key []byte) error {
	// concatenating card to string
	cardStr := card.Cardname + "," + card.Number + "," + card.Name + "," + card.Surname + "," + card.ValidTill + "," + card.Code

	// encrypting the card data
	encrypted, err := encrypt.EncryptData([]byte(cardStr), key)
	if err != nil {
		return err
	}
	// Encode the encrypted password in base64 for storage
	encodedData := base64.StdEncoding.EncodeToString(encrypted)

	file, err := os.OpenFile(config.ClientCfg.LocalStorage+s.Cards, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("error in SaveCard when opening file:%w", err)
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "%s\n", encodedData)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return fmt.Errorf("error in SaveCard when writing in file:%w", err)
	}

	return nil
}

func (s *LocalStorage) GetCard(cardname string, key []byte) (card database.Card, err error) {

	// opening the localstorage file
	file, err := os.OpenFile(config.ClientCfg.LocalStorage+s.Cards, os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		return database.Card{}, fmt.Errorf("error in SaveCard when opening file:%w", err)
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

func (s *LocalStorage) ListCards(key []byte) (cards []string, err error) {

	// opening the localstorage file
	file, err := os.OpenFile(config.ClientCfg.LocalStorage+s.Cards, os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		return nil, fmt.Errorf("error in ListCards when opening file:%w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		scanner.Text()
	}

	return cards, nil
}
