package localstorage

import (
	"fmt"
	"os"

	"github.com/gambruh/gophkeeper/internal/config"
	"github.com/gambruh/gophkeeper/internal/database"
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

func (s *LocalStorage) SaveCard(card database.Card) error {

	file, err := os.OpenFile(config.ClientCfg.LocalStorage+s.Cards, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("error in SaveCard when opening file:%w", err)
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "%s,%s,%s,%s,%s,%s\n", card.Cardname, card.Number, card.Name, card.Surname, card.ValidTill, card.Code)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return fmt.Errorf("error in SaveCard when writing in file:%w", err)
	}

	return nil
}
