package clientfunc

import (
	"context"
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
			err := c.checkCards()
			if err != nil {
				log.Println("error in DataChecker function returned from CheckCards:", err)
			}
		}
	}
}

func (c *Client) CheckAll() error {

	err := c.checkCards()
	if err != nil {
		return err
	}

	err = c.checkLoginCreds()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) checkCards() error {
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

	mapLocal := helpers.CreateMapFromList(cardsLocal)
	mapServer := helpers.CreateMapFromList(cardsDB)

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

func (c *Client) checkLoginCreds() error {
	if c.AuthCookie == nil {
		return nil
	}

	loginCredsDB, err := c.listLoginCredsFromDB()
	if err != nil {
		return err
	}

	loginCredsLocal, err := c.listLoginCredsFromStorage()
	if err != nil {
		return err
	}

	mapLocal := helpers.CreateMapFromList(loginCredsLocal)
	mapServer := helpers.CreateMapFromList(loginCredsDB)

	toUpload, toDownload := helpers.CompareTwoMaps(mapServer, mapLocal)

	for logincredname := range toUpload {
		var eLoginCred database.EncryptedData
		logincred, err := c.getLoginCredsFromStorage(logincredname)
		if err != nil {
			return err
		}
		eLoginCred.Data, err = helpers.EncryptLoginCredsData(logincred, c.Key)
		if err != nil {
			return err
		}
		eLoginCred.Name = logincred.Name
		err = c.sendLoginCredsToDB(eLoginCred)
		if err != nil {
			return err
		}
	}

	for logincredname := range toDownload {
		eLoginCred, err := c.getLoginCredsFromDB(logincredname)

		if err != nil {
			return err
		}

		logincred, err := helpers.DecryptLoginCredsData(eLoginCred, c.Key)
		if err != nil {
			return err
		}

		err = c.saveLoginCredsInStorage(logincred)
		if err != nil {
			return err
		}
	}

	return nil
}
