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
	//fmt.Println("cards from storage:", cards)
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
