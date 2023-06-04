package clientfunc

import (
	"bytes"
	"encoding/json"
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
	case 500:
		fmt.Println("server is down")
	}
	return nil
}

func (c *Client) saveCardInStorage(card database.Card) error {
	err := c.Storage.SaveCard(card)
	if err != nil {
		return fmt.Errorf("error in saveCardInStorage:%w", err)
	}
	return nil
}
