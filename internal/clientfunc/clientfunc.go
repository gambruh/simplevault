package clientfunc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gambruh/gophkeeper/internal/auth"
	"github.com/gambruh/gophkeeper/internal/config"
	"github.com/gambruh/gophkeeper/internal/database"
)

type Client struct {
	Server       string
	LocalStorage database.Storage
	Client       *http.Client
	AuthCookie   *http.Cookie
}

func NewClient() *Client {

	return &Client{
		Server:       config.ClientCfg.Address,
		LocalStorage: database.NewStorage(),
		Client:       &http.Client{},
	}

}

// PrintAvailableCommands function prints available commands
func PrintAvailableCommands(commands map[string]func([]string)) {
	fmt.Println("Available commands:")
	for cmd := range commands {
		fmt.Println("-", cmd)
	}
}

func (c *Client) Register(input []string) {
	if len(input) != 3 {
		printRegisterSyntacsys()
		return
	}
	var loginData auth.LoginData
	for i, data := range input {
		switch i {
		case 0:
		case 1:
			loginData.Login = data
		case 2:
			loginData.Password = data
		}
	}

	authcookie, err := c.sendRegisterRequest(loginData)

	switch err {
	case nil:
		c.AuthCookie = authcookie
		fmt.Println("registered successfully")
	case ErrWrongLoginData:
		fmt.Println("please provide correct login data")
	default:
		fmt.Println("encountered error when trying to register new user: ", err)
	}
}

func (c *Client) Login(input []string) {
	if len(input) != 3 {
		printLoginSyntacsys()
		return
	}
	var loginData auth.LoginData
	for i, data := range input {
		switch i {
		case 0:
		case 1:
			loginData.Login = data
		case 2:
			loginData.Password = data
		}
	}

	authcookie, err := c.sendLoginRequest(loginData)
	if err != nil {
		fmt.Println("encountered error when trying to register new user: ", err)
		return
	}
	c.AuthCookie = authcookie
}

func (c *Client) SetCardCommand(input []string) {
	if c.AuthCookie == nil {
		fmt.Println("please login first")
		return
	}
	if len(input) != 7 {
		printSetCardSyntacsys()
		return
	}
	var cardData database.Card
	for i, data := range input {
		switch i {
		case 0:
		case 1:
			cardData.Bank = data
		case 2:
			cardData.Number = data
		case 3:
			cardData.Name = data
		case 4:
			cardData.Surname = data
		case 5:
			cardData.ValidTill = data
		case 6:
			cardData.Code = data
		}
	}

	err := c.sendCardToStorage(cardData)
	if err != nil {
		fmt.Println("error in client sending data to database in SetCardCommand:", err)
	}
}

func (c *Client) sendCardToStorage(card database.Card) error {
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

func (c *Client) sendRegisterRequest(login auth.LoginData) (*http.Cookie, error) {

	url := fmt.Sprintf("%s/api/user/register", c.Server)

	if !strings.HasPrefix(url, "http://") {
		url = "http://" + url
	}

	jsbody, err := json.Marshal(login)
	if err != nil {
		return nil, fmt.Errorf("error when marshaling json: %w", err)
	}
	rbody := bytes.NewBuffer(jsbody)
	r, err := http.NewRequest(http.MethodPost, url, rbody)
	if err != nil {
		return nil, fmt.Errorf("error when creating NewRequest: %w", err)
	}
	r.Header.Add("Content-Type", "application/json")
	res, err := c.Client.Do(r)
	if err != nil {
		return nil, fmt.Errorf("error when sending request in sendLoginRequest: %w", err)
	}
	defer res.Body.Close()

	cookies := res.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "gophkeeper-auth" {
			return cookie, nil
		}
	}
	return nil, ErrNoCookieReturned
}

func (c *Client) sendLoginRequest(login auth.LoginData) (*http.Cookie, error) {

	url := fmt.Sprintf("%s/api/user/login", c.Server)

	if !strings.HasPrefix(url, "http://") {
		url = "http://" + url
	}

	jsbody, err := json.Marshal(login)
	if err != nil {
		return nil, fmt.Errorf("error when marshaling json: %w", err)
	}
	rbody := bytes.NewBuffer(jsbody)
	r, err := http.NewRequest(http.MethodPost, url, rbody)
	if err != nil {
		return nil, fmt.Errorf("error when creating NewRequest: %w", err)
	}
	r.Header.Add("Content-Type", "application/json")
	res, err := c.Client.Do(r)
	if err != nil {
		return nil, fmt.Errorf("error when sending request in sendLoginRequest: %w", err)
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case 200:
		cookies := res.Cookies()
		for _, cookie := range cookies {
			if cookie.Name == "gophkeeper-auth" {
				return cookie, nil
			}
		}
		return nil, ErrNoCookieReturned
	case 401:
		return nil, ErrWrongLoginData
	case 500:
		return nil, ErrServerIsDown
	default:
		return nil, fmt.Errorf("unexpected error in sendLoginRequest")
	}

	//return nil, fmt.Errorf("error when sending request in sendLoginRequest: %w", err)
}

func printSetCardSyntacsys() {
	fmt.Println("Wrong input!")
	fmt.Println("Right syntacsis: setcard <cardname> <cardnumber> <cardholder name> <cardholder surname> <card valid till date in format 'dd:mm:yyyy'> <cvv code>")
}

func printRegisterSyntacsys() {
	fmt.Println("Wrong input!")
	fmt.Println("Right syntacsis: register <login> <password>")
}

func printLoginSyntacsys() {
	fmt.Println("Wrong input!")
	fmt.Println("Right syntacsis: login <login> <password>")
}
