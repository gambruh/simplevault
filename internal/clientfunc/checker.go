package clientfunc

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gambruh/gophkeeper/internal/database"
	"github.com/gambruh/gophkeeper/internal/helpers"
)

// DataChecker synchronizes data between local storage and DB
func (c *Client) DataChecker(context context.Context, ticker *time.Ticker, wgShutdown *sync.WaitGroup) {
	defer wgShutdown.Done()

	for {
		select {
		case <-context.Done():
			return
		case <-ticker.C:
			err := c.CheckCards()
			if err != nil {
				log.Println("error in CheckAll function returned from CheckCards:", err)
			}
			fmt.Println("vault synced!")
		}
	}

}

func (c *Client) CheckCards() error {
	if c.AuthCookie == nil {
		return nil
	}
	cardsDB, err := c.listCardsFromDB()
	if err != nil {
		return err
	}

	cardsLocal, err := c.listCardsFromStorage()
	if err != nil {
		return err
	}

	mapLocal := make(map[string]struct{})
	mapServer := make(map[string]struct{})

	for _, cardname := range cardsLocal {
		mapLocal[cardname] = struct{}{}
	}

	for _, cardname := range cardsDB {
		mapServer[cardname] = struct{}{}
	}

	toUpload, toDownload := helpers.CompareTwoMaps(mapServer, mapLocal)

	for cardname := range toUpload {
		var eCard database.EncryptedCard
		card, err := c.getCardFromStorage(cardname)
		if err != nil {
			return err
		}
		eCard.Data, err = helpers.EncryptCardData(card, c.Key)
		if err != nil {
			return err
		}
		eCard.Cardname = card.Cardname
		err = c.sendCardToDB(eCard)
		if err != nil {
			return err
		}
	}

	for cardname := range toDownload {
		eCard, err := c.getCardFromDB(cardname)

		if err != nil {
			return err
		}

		card, err := helpers.DecryptCardData(eCard, c.Key)
		if err != nil {
			return err
		}

		err = c.saveCardInStorage(card)
		if err != nil {
			return err
		}
	}

	return nil
}
