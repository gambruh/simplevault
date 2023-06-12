package clientfunc

import (
	"fmt"

	"github.com/gambruh/gophkeeper/internal/database"
	"github.com/gambruh/gophkeeper/internal/helpers"
)

// CheckAll function checks for desynchronized data between client and server
func (c *Client) CheckAll() error {

	if err := c.checkCards(); err != nil {
		if err == ErrDataNotFound {
			//do nothing
		} else {
			return fmt.Errorf("error in checkCards:%w", err)
		}
	}

	if err := c.checkLoginCreds(); err != nil {
		if err == ErrDataNotFound {
			//do nothing
		} else {
			return fmt.Errorf("error in checkLoginCreds:%w", err)
		}
	}

	if err := c.checkNotes(); err != nil {
		if err == ErrDataNotFound {
			//do nothing
		} else {
			return fmt.Errorf("error in checkNotes:%w", err)
		}
	}

	if err := c.checkBinaries(); err != nil {
		if err == ErrDataNotFound {
			//do nothing
		} else {
			return fmt.Errorf("error in checkBinaries:%w", err)
		}

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

func (c *Client) checkNotes() error {
	if c.AuthCookie == nil {
		return nil
	}

	listDB, err := c.listNotesFromDB()
	if err != nil {
		return err
	}

	listLocal, err := c.listNotesFromStorage()
	if err != nil {
		return err
	}

	mapLocal := helpers.CreateMapFromList(listLocal)
	mapServer := helpers.CreateMapFromList(listDB)

	toUpload, toDownload := helpers.CompareTwoMaps(mapServer, mapLocal)

	for notename := range toUpload {
		var eData database.EncryptedData
		note, err := c.getNoteFromStorage(notename)
		if err != nil {
			return err
		}
		eData.Data, err = helpers.EncryptNoteData(note, c.Key)
		if err != nil {
			return err
		}
		eData.Name = note.Name
		err = c.sendNoteToDB(eData)
		if err != nil {
			return err
		}
	}

	for notename := range toDownload {
		eData, err := c.getNoteFromDB(notename)

		if err != nil {
			return err
		}

		note, err := helpers.DecryptNoteData(eData, c.Key)
		if err != nil {
			return err
		}

		err = c.saveNoteInStorage(note)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) checkBinaries() error {
	if c.AuthCookie == nil {
		return nil
	}

	listDB, err := c.listBinariesFromDB()
	if err != nil {
		return err
	}

	listLocal, err := c.listBinariesFromStorage()
	if err != nil {
		return err
	}

	mapLocal := helpers.CreateMapFromList(listLocal)
	mapServer := helpers.CreateMapFromList(listDB)

	toUpload, toDownload := helpers.CompareTwoMaps(mapServer, mapLocal)

	for binaryname := range toUpload {
		binary, err := c.getBinaryFromStorage(binaryname)
		if err != nil {
			return err
		}

		err = c.sendBinaryToDB(binary)
		if err != nil {
			return err
		}
	}

	for binaryname := range toDownload {
		binary, err := c.getBinaryFromDB(binaryname)
		if err != nil {
			return err
		}

		err = c.saveBinaryInStorage(binary)
		if err != nil {
			return err
		}
	}

	return nil
}
