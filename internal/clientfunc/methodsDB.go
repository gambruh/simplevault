package clientfunc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gambruh/gophkeeper/internal/storage"
)

func (c *Client) sendCardToDB(encrCard storage.EncryptedData) error {
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

func (c *Client) getCardFromDB(cardname string) (card storage.EncryptedData, err error) {
	var inCard storage.EncryptedData
	inCard.Name = cardname
	url := fmt.Sprintf("%s/api/cards/get", c.Server)
	if !strings.HasPrefix(url, "http://") {
		url = "http://" + url
	}

	jsbody, err := json.Marshal(inCard)
	if err != nil {
		return storage.EncryptedData{}, fmt.Errorf("error when marshaling json in getCardFromDB: %w", err)
	}
	rbody := bytes.NewBuffer(jsbody)
	r, err := http.NewRequest(http.MethodPost, url, rbody)
	if err != nil {
		return storage.EncryptedData{}, fmt.Errorf("error when creating NewRequest in getCardFromDB: %w", err)
	}
	r.Header.Add("Content-Type", "application/json")
	r.AddCookie(c.AuthCookie)
	res, err := c.Client.Do(r)
	if err != nil {
		return storage.EncryptedData{}, fmt.Errorf("error when sending request in getCardFromDB: %w", err)
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case 200:
		err := json.NewDecoder(res.Body).Decode(&card)
		if err != nil {
			return storage.EncryptedData{}, fmt.Errorf("error when decoding json in getCardFromDB: %w", err)
		}
		return card, nil
	case 204:
		return storage.EncryptedData{}, ErrDataNotFound
	case 400:
		return storage.EncryptedData{}, ErrBadRequest
	case 401:
		return storage.EncryptedData{}, ErrLoginRequired
	case 500:
		return storage.EncryptedData{}, ErrServerIsDown
	default:
		return storage.EncryptedData{}, errors.New("unexpected error")
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

func (c *Client) SendCardToDB(cardData storage.EncryptedData) {
	err := c.sendCardToDB(cardData)
	switch err {
	case ErrMetanameIsTaken:
		log.Println("There are already card with this name in database. Please provide new cardname or edit current")
	case ErrLoginRequired:
		log.Println("Please login to the server")
	case nil:
		log.Printf("Saved card %s to the vault\n", cardData.Name)
	default:
		log.Println("error in client sending data to database in SetCardCommand:", err)
	}
}

func (c *Client) sendLoginCredsToDB(encrData storage.EncryptedData) error {
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

func (c *Client) getLoginCredsFromDB(logincredname string) (encrData storage.EncryptedData, err error) {
	var input storage.EncryptedData
	input.Name = logincredname
	url := fmt.Sprintf("%s/api/logincreds/get", c.Server)
	if !strings.HasPrefix(url, "http://") {
		url = "http://" + url
	}

	jsbody, err := json.Marshal(input)
	if err != nil {
		return storage.EncryptedData{}, fmt.Errorf("error when marshaling json in getLoginCredsFromDB: %w", err)
	}
	rbody := bytes.NewBuffer(jsbody)
	r, err := http.NewRequest(http.MethodPost, url, rbody)
	if err != nil {
		return storage.EncryptedData{}, fmt.Errorf("error when creating NewRequest in getLoginCredsFromDB: %w", err)
	}
	r.Header.Add("Content-Type", "application/json")
	r.AddCookie(c.AuthCookie)
	res, err := c.Client.Do(r)
	if err != nil {
		return storage.EncryptedData{}, fmt.Errorf("error when sending request in getLoginCredsFromDB: %w", err)
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case 200:
		err := json.NewDecoder(res.Body).Decode(&encrData)
		if err != nil {
			return storage.EncryptedData{}, fmt.Errorf("error when decoding json in getLoginCredsFromDB: %w", err)
		}
		return encrData, nil
	case 204:
		return storage.EncryptedData{}, ErrDataNotFound
	case 400:
		return storage.EncryptedData{}, ErrBadRequest
	case 401:
		return storage.EncryptedData{}, ErrLoginRequired
	case 500:
		return storage.EncryptedData{}, ErrServerIsDown
	default:
		return storage.EncryptedData{}, errors.New("unexpected error")
	}
}

func (c *Client) sendNoteToDB(encrData storage.EncryptedData) error {
	url := fmt.Sprintf("%s/api/notes/add", c.Server)

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

func (c *Client) listNotesFromDB() (notes []string, err error) {
	url := fmt.Sprintf("%s/api/notes/list", c.Server)
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
		err := json.NewDecoder(res.Body).Decode(&notes)
		if err != nil {
			return nil, fmt.Errorf("error when decoding json in listCardsFromDB: %w", err)
		}
	case 401:
		return nil, ErrLoginRequired
	case 500:
		return nil, ErrServerIsDown
	}
	return notes, nil
}

func (c *Client) getNoteFromDB(notename string) (encrData storage.EncryptedData, err error) {
	var input storage.EncryptedData
	input.Name = notename
	url := fmt.Sprintf("%s/api/notes/get", c.Server)
	if !strings.HasPrefix(url, "http://") {
		url = "http://" + url
	}

	jsbody, err := json.Marshal(input)
	if err != nil {
		return storage.EncryptedData{}, fmt.Errorf("error when marshaling json in getNoteFromDB: %w", err)
	}
	rbody := bytes.NewBuffer(jsbody)
	r, err := http.NewRequest(http.MethodPost, url, rbody)
	if err != nil {
		return storage.EncryptedData{}, fmt.Errorf("error when creating NewRequest in getNoteFromDB: %w", err)
	}
	r.Header.Add("Content-Type", "application/json")
	r.AddCookie(c.AuthCookie)
	res, err := c.Client.Do(r)
	if err != nil {
		return storage.EncryptedData{}, fmt.Errorf("error when sending request in getNoteFromDB: %w", err)
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case 200:
		err := json.NewDecoder(res.Body).Decode(&encrData)
		if err != nil {
			return storage.EncryptedData{}, fmt.Errorf("error when decoding json in getNoteFromDB: %w", err)
		}
		return encrData, nil
	case 204:
		return storage.EncryptedData{}, ErrDataNotFound
	case 400:
		return storage.EncryptedData{}, ErrBadRequest
	case 401:
		return storage.EncryptedData{}, ErrLoginRequired
	case 500:
		return storage.EncryptedData{}, ErrServerIsDown
	default:
		return storage.EncryptedData{}, errors.New("unexpected error")
	}
}

func (c *Client) sendBinaryToDB(encrData storage.Binary) error {
	url := fmt.Sprintf("%s/api/binaries/add", c.Server)

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

func (c *Client) getBinaryFromDB(binaryname string) (encrData storage.Binary, err error) {
	var input storage.Binary
	input.Name = binaryname
	url := fmt.Sprintf("%s/api/binaries/get", c.Server)
	if !strings.HasPrefix(url, "http://") {
		url = "http://" + url
	}

	jsbody, err := json.Marshal(input)
	if err != nil {
		return storage.Binary{}, fmt.Errorf("error when marshaling json in getBinaryFromDB: %w", err)
	}
	rbody := bytes.NewBuffer(jsbody)
	r, err := http.NewRequest(http.MethodPost, url, rbody)
	if err != nil {
		return storage.Binary{}, fmt.Errorf("error when creating NewRequest in getBinaryFromDB: %w", err)
	}
	r.Header.Add("Content-Type", "application/json")
	r.AddCookie(c.AuthCookie)
	res, err := c.Client.Do(r)
	if err != nil {
		return storage.Binary{}, fmt.Errorf("error when sending request in getBinaryFromDB: %w", err)
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case 200:
		err := json.NewDecoder(res.Body).Decode(&encrData)
		if err != nil {
			return storage.Binary{}, fmt.Errorf("error when decoding json in getBinaryFromDB: %w", err)
		}
		return encrData, nil
	case 204:
		return storage.Binary{}, ErrDataNotFound
	case 400:
		return storage.Binary{}, ErrBadRequest
	case 401:
		return storage.Binary{}, ErrLoginRequired
	case 500:
		return storage.Binary{}, ErrServerIsDown
	default:
		return storage.Binary{}, errors.New("unexpected error")
	}
}

func (c *Client) listBinariesFromDB() (binaries []string, err error) {
	url := fmt.Sprintf("%s/api/binaries/list", c.Server)
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
		err := json.NewDecoder(res.Body).Decode(&binaries)
		if err != nil {
			return nil, fmt.Errorf("error when decoding json in listCardsFromDB: %w", err)
		}
	case 401:
		return nil, ErrLoginRequired
	case 500:
		return nil, ErrServerIsDown
	}
	return binaries, nil
}
