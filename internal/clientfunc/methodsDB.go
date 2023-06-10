package clientfunc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gambruh/gophkeeper/internal/database"
)

func (c *Client) sendCardToDB(encrCard database.EncryptedCard) error {
	url := fmt.Sprintf("%s/api/cards/add", c.Server)

	if !strings.HasPrefix(url, "http://") {
		url = "http://" + url
	}

	jsbody, err := json.Marshal(encrCard)
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

func (c *Client) getCardFromDB(cardname string) (card database.EncryptedCard, err error) {
	var inCard database.EncryptedCard
	inCard.Cardname = cardname
	url := fmt.Sprintf("%s/api/cards/get", c.Server)
	if !strings.HasPrefix(url, "http://") {
		url = "http://" + url
	}

	jsbody, err := json.Marshal(inCard)
	if err != nil {
		return database.EncryptedCard{}, fmt.Errorf("error when marshaling json in getCardFromDB: %w", err)
	}
	rbody := bytes.NewBuffer(jsbody)
	r, err := http.NewRequest(http.MethodPost, url, rbody)
	if err != nil {
		return database.EncryptedCard{}, fmt.Errorf("error when creating NewRequest in getCardFromDB: %w", err)
	}
	r.Header.Add("Content-Type", "application/json")
	r.AddCookie(c.AuthCookie)
	res, err := c.Client.Do(r)
	if err != nil {
		return database.EncryptedCard{}, fmt.Errorf("error when sending request in getCardFromDB: %w", err)
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case 200:
		err := json.NewDecoder(res.Body).Decode(&card)
		if err != nil {
			return database.EncryptedCard{}, fmt.Errorf("error when decoding json in getCardFromDB: %w", err)
		}
		return card, nil
	case 204:
		return database.EncryptedCard{}, ErrDataNotFound
	case 400:
		return database.EncryptedCard{}, ErrBadRequest
	case 401:
		return database.EncryptedCard{}, ErrLoginRequired
	case 500:
		return database.EncryptedCard{}, ErrServerIsDown
	default:
		return database.EncryptedCard{}, errors.New("unexpected error")
	}
}

func (c *Client) GetCardFromDB(cardname string) {
	card, err := c.getCardFromDB(cardname)
	switch err {
	case ErrDataNotFound:
		fmt.Println("Card with that name not found in the DB")
	case ErrLoginRequired:
		fmt.Println("Please login to the server")
	case ErrBadRequest:
		fmt.Println("Please contact devs to change API interaction, wrong request")
	case ErrServerIsDown:
		fmt.Println("Internal server error")
	case nil:
		fmt.Printf("%+v\n", card)
	}
}

func (c *Client) SendCardToDB(cardData database.EncryptedCard) {
	err := c.sendCardToDB(cardData)
	switch err {
	case ErrMetanameIsTaken:
		log.Println("There are already card with this name in database. Please provide new cardname or edit current")
	case ErrLoginRequired:
		log.Println("Please login to the server")
	case nil:
		log.Printf("Saved card %s to the vault\n", cardData.Cardname)
	default:
		log.Println("error in client sending data to database in SetCardCommand:", err)
	}
}

func (c *Client) sendLoginCredsToDB(encrData database.EncryptedData) error {
	url := fmt.Sprintf("%s/api/logincreds/add", c.Server)

	if !strings.HasPrefix(url, "http://") {
		url = "http://" + url
	}

	jsbody, err := json.Marshal(encrData)
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
		return fmt.Errorf("error when sending request in sendLoginCredsToDB: %w", err)
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

func (c *Client) listLoginCredsFromDB() (logincreds []string, err error) {
	url := fmt.Sprintf("%s/api/logincreds/list", c.Server)
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
		err := json.NewDecoder(res.Body).Decode(&logincreds)
		if err != nil {
			return nil, fmt.Errorf("error when decoding json in listCardsFromDB: %w", err)
		}
	case 401:
		return nil, ErrLoginRequired
	case 500:
		return nil, ErrServerIsDown
	}
	return logincreds, nil
}

func (c *Client) getLoginCredsFromDB(logincredname string) (encrData database.EncryptedData, err error) {
	var input database.EncryptedData
	input.Name = logincredname
	url := fmt.Sprintf("%s/api/logincreds/get", c.Server)
	if !strings.HasPrefix(url, "http://") {
		url = "http://" + url
	}

	jsbody, err := json.Marshal(input)
	if err != nil {
		return database.EncryptedData{}, fmt.Errorf("error when marshaling json in getLoginCredsFromDB: %w", err)
	}
	rbody := bytes.NewBuffer(jsbody)
	r, err := http.NewRequest(http.MethodPost, url, rbody)
	if err != nil {
		return database.EncryptedData{}, fmt.Errorf("error when creating NewRequest in getLoginCredsFromDB: %w", err)
	}
	r.Header.Add("Content-Type", "application/json")
	r.AddCookie(c.AuthCookie)
	res, err := c.Client.Do(r)
	if err != nil {
		return database.EncryptedData{}, fmt.Errorf("error when sending request in getLoginCredsFromDB: %w", err)
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case 200:
		err := json.NewDecoder(res.Body).Decode(&encrData)
		if err != nil {
			return database.EncryptedData{}, fmt.Errorf("error when decoding json in getLoginCredsFromDB: %w", err)
		}
		return encrData, nil
	case 204:
		return database.EncryptedData{}, ErrDataNotFound
	case 400:
		return database.EncryptedData{}, ErrBadRequest
	case 401:
		return database.EncryptedData{}, ErrLoginRequired
	case 500:
		return database.EncryptedData{}, ErrServerIsDown
	default:
		return database.EncryptedData{}, errors.New("unexpected error")
	}
}
