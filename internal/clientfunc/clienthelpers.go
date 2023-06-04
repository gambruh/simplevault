package clientfunc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gambruh/gophkeeper/internal/database"
)

// Card commands helpers

func (c *Client) sendCardToDB(card database.Card) error {
	url := fmt.Sprintf("%s/api/cards/add", c.Server)

	if !strings.HasPrefix(url, "http://") {
		url = "http://" + url
	}

	jsbody, err := json.Marshal(card)
	if err != nil {
		return fmt.Errorf("error when marshaling json: %w", err)
	}
	rbody := bytes.NewBuffer(jsbody)
	r, err := http.NewRequest(http.MethodPost, url, rbody)
	if err != nil {
		return fmt.Errorf("error when creating NewRequest: %w", err)
	}
	r.Header.Add("Content-Type", "application/json")
	r.AddCookie(c.AuthCookie)
	res, err := c.Client.Do(r)
	if err != nil {
		return fmt.Errorf("error when sending request in sendCardToStorage: %w", err)
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case 200:
		return nil
	case 401:
		return ErrLoginRequired
	case 409:
		return ErrMetanameIsTaken
	case 500:
		fmt.Println("internal server error")
	}
	return nil
}

func (c *Client) saveCardInStorage(card database.Card) error {
	err := c.Storage.SaveCard(card, c.Key)
	if err != nil {
		return fmt.Errorf("error in saveCardInStorage:%w", err)
	}
	return nil
}

func (c *Client) listCardsFromDB() (cards []string, err error) {
	url := fmt.Sprintf("%s/api/cards/list", c.Server)
	if !strings.HasPrefix(url, "http://") {
		url = "http://" + url
	}

	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error when creating NewRequest: %w", err)
	}
	r.Header.Add("Content-Type", "application/json")
	r.AddCookie(c.AuthCookie)
	res, err := c.Client.Do(r)
	if err != nil {
		return nil, fmt.Errorf("error when sending request in listCardsFromDB: %w", err)
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case 200:
		err := json.NewDecoder(res.Body).Decode(&cards)
		if err != nil {
			return nil, fmt.Errorf("error when decoding json in listCardsFromDB: %w", err)
		}
	case 401:
		return nil, ErrLoginRequired
	case 500:
		return nil, ErrServerIsDown
	}

	return cards, nil
}

func (c *Client) getCardFromDB(cardname string) (card database.Card, err error) {
	var inCard database.Card
	inCard.Cardname = cardname
	url := fmt.Sprintf("%s/api/cards/get", c.Server)
	if !strings.HasPrefix(url, "http://") {
		url = "http://" + url
	}

	jsbody, err := json.Marshal(inCard)
	if err != nil {
		return database.Card{}, fmt.Errorf("error when marshaling json in getCardFromDB: %w", err)
	}
	rbody := bytes.NewBuffer(jsbody)
	r, err := http.NewRequest(http.MethodPost, url, rbody)
	if err != nil {
		return database.Card{}, fmt.Errorf("error when creating NewRequest in getCardFromDB: %w", err)
	}
	r.Header.Add("Content-Type", "application/json")
	r.AddCookie(c.AuthCookie)
	res, err := c.Client.Do(r)
	if err != nil {
		return database.Card{}, fmt.Errorf("error when sending request in getCardFromDB: %w", err)
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case 200:
		err := json.NewDecoder(res.Body).Decode(&card)
		if err != nil {
			return database.Card{}, fmt.Errorf("error when decoding json in getCardFromDB: %w", err)
		}
		return card, nil
	case 204:
		return database.Card{}, ErrDataNotFound
	case 400:
		return database.Card{}, ErrBadRequest
	case 401:
		return database.Card{}, ErrLoginRequired
	case 500:
		return database.Card{}, ErrServerIsDown
	default:
		return database.Card{}, errors.New("unexpected error")
	}

}

func (c *Client) getCardFromLocalStorage(cardname string) (card database.Card, err error) {

	card, err = c.Storage.GetCard(cardname, c.Key)
	if err != nil {
		return database.Card{}, err
	}

	return card, nil
}
