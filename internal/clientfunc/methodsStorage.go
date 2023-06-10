package clientfunc

import (
	"fmt"
	"log"

	"github.com/gambruh/gophkeeper/internal/database"
)

// Card commands helpers

func (c *Client) saveCardInStorage(card database.Card) error {
	err := c.Storage.SaveCard(card, c.Key)
	if err != nil {
		return fmt.Errorf("error in saveCardInStorage:%w", err)
	}
	return nil
}

func (c *Client) listCardsFromStorage() (cards []string, err error) {

	cards, err = c.Storage.ListCards()
	if err != nil {
		return nil, err
	}

	return cards, nil
}

func (c *Client) getCardFromStorage(cardname string) (card database.Card, err error) {

	card, err = c.Storage.GetCard(cardname, c.Key)
	if err != nil {
		fmt.Println("err in getCardFromStorage is:", err)
		return database.Card{}, err
	}

	return card, nil
}

func (c *Client) DeleteLocalStorage() {

	err := c.Storage.DeleteLocalStorage()

	if err != nil {
		log.Println(err)
	}

}

func (c *Client) saveLoginCredsInStorage(logincreds database.LoginCreds) error {

	err := c.Storage.SaveLoginCreds(logincreds, c.Key)
	if err != nil {
		return fmt.Errorf("error in saveCardInStorage:%w", err)
	}
	return nil

}

func (c *Client) getLoginCredsFromStorage(logincredname string) (logincred database.LoginCreds, err error) {
	logincred, err = c.Storage.GetLoginCreds(logincredname, c.Key)
	if err != nil {
		fmt.Println("err in getLoginCredsFromStorage is:", err)
		return database.LoginCreds{}, err
	}

	return logincred, nil
}

func (c *Client) listLoginCredsFromStorage() (logincreds []string, err error) {
	logincreds, err = c.Storage.ListLoginCreds()
	if err != nil {
		return nil, err
	}

	return logincreds, nil
}

func (c *Client) saveNoteInStorage(note database.Note) error {

	err := c.Storage.SaveNote(note, c.Key)
	if err != nil {
		return fmt.Errorf("error in saveCardInStorage:%w", err)
	}
	return nil

}

func (c *Client) getNoteFromStorage(notename string) (note database.Note, err error) {
	note, err = c.Storage.GetNote(notename, c.Key)
	if err != nil {
		fmt.Println("err in getNoteFromStorage is:", err)
		return database.Note{}, err
	}

	return note, nil
}

func (c *Client) listNotesFromStorage() (notes []string, err error) {
	notes, err = c.Storage.ListNotes()
	if err != nil {
		return nil, err
	}

	return notes, nil
}

func (c *Client) saveBinaryInStorage(binary database.Binary) error {
	err := c.Storage.SaveBinary(binary, c.Key)
	if err != nil {
		fmt.Println("couldn't save binary into storage:", err)
		return err
	}

	return nil
}

func (c *Client) getBinaryFromStorage(binaryname string) (binary database.Binary, err error) {
	binary, err = c.Storage.GetBinary(binaryname, c.Key)
	if err != nil {
		fmt.Println("err in getBinaryFromStorage :", err)
		return database.Binary{}, err
	}

	return binary, nil
}

func (c *Client) listBinariesFromStorage() (binaries []string, err error) {
	binaries, err = c.Storage.ListBinaries()
	if err != nil {
		return nil, err
	}

	return binaries, nil
}
