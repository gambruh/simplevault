package clientfunc

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gambruh/gophkeeper/internal/argon2id"
	"github.com/gambruh/gophkeeper/internal/auth"
	"github.com/gambruh/gophkeeper/internal/config"
)

const userdata = "./userdata/user.json"

func getUserDataFromFile() (auth.LoginData, error) {

	var logindata auth.LoginData
	ufile, err := os.Open(userdata)
	if err != nil {
		fmt.Println("please login, using login command")
		return auth.LoginData{}, err
	}

	err = json.NewDecoder(ufile).Decode(&logindata)
	if err != nil {
		fmt.Println("please delete user.json and relogin, using login command - cant unmarshal the file")
		return auth.LoginData{}, err
	}
	defer ufile.Close()

	return logindata, nil
}

func (c *Client) Register(input []string) {
	if len(input) != 3 {
		printRegisterSyntax()
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
		key := sha256.Sum256([]byte(loginData.Password))
		c.Key = key[:]
		fmt.Println("registered successfully")
		c.CreateUserLoginFile(loginData.Login, loginData.Password)
	case ErrUsernameIsTaken:
		fmt.Println("Username is taken, please provide another")
	default:
		fmt.Println("error when trying to register new user: ", err)
	}
}

func (c *Client) Login(input []string) {
	if len(input) != 3 {
		printLoginSyntax()
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

	// logging into server
	authcookie, err := c.sendLoginRequest(loginData)
	switch err {
	case nil:
		fmt.Println("Logged on the server!")
		c.AuthCookie = authcookie
		key := sha256.Sum256([]byte(loginData.Password))
		c.Key = key[:]
		c.CreateUserLoginFile(loginData.Login, loginData.Password)
	case ErrWrongLoginData:
		fmt.Println("wrong login credentials, try again")
		return
	case ErrServerIsDown:
		fmt.Println("Server is down, try again later")
		c.AuthCookie = nil
	default:
		fmt.Println("error when trying to login online: ", err)
	}

	// logging offline
	err = c.loginOffline(loginData)
	if err != nil {
		fmt.Println("error when trying to login offline: ", err)
		return
	}
	fmt.Println("Successfully logged offline!")
}

func (c *Client) loginOffline(logincreds auth.LoginData) error {

	checklogindata, err := getUserDataFromFile()
	if err != nil {
		return err
	}

	if logincreds.Login != checklogindata.Login {
		return ErrWrongLoginData
	}

	hashCheck, err := argon2id.ComparePasswordAndHash(logincreds.Password, checklogindata.Password)
	if err != nil {
		return fmt.Errorf("error when trying to compare hashes:%w", err)
	}

	if !hashCheck {
		return ErrWrongLoginData
	}

	//sucessfuly logged in
	c.LoggedOffline = true
	return nil
}

func (c *Client) sendRegisterRequest(login auth.LoginData) (*http.Cookie, error) {
	//preparing url to send to
	url := fmt.Sprintf("%s/api/user/register", c.Server)
	//checking if the prefix is ok
	if !strings.HasPrefix(url, "http://") {
		url = "http://" + url
	}

	//preparing json body
	jsbody, err := json.Marshal(login)
	if err != nil {
		return nil, fmt.Errorf("error when marshaling json: %w", err)
	}
	rbody := bytes.NewBuffer(jsbody)

	//preparing request
	r, err := http.NewRequest(http.MethodPost, url, rbody)
	if err != nil {
		return nil, fmt.Errorf("error when creating NewRequest: %w", err)
	}
	r.Header.Add("Content-Type", "application/json")

	//sending request
	res, err := c.Client.Do(r)
	if err != nil {
		return nil, fmt.Errorf("error in sendLoginRequest: %w", err)
	}
	defer res.Body.Close()

	// checking the response code
	switch res.StatusCode {
	case 200:
		cookies := res.Cookies()
		for _, cookie := range cookies {
			if cookie.Name == "gophkeeper-auth" {
				return cookie, nil
			}
		}
		return nil, ErrNoCookieReturned
	case 409:
		return nil, ErrUsernameIsTaken
	case 500:
		fmt.Println("Server error, please try again")
		return nil, ErrServerIsDown
	}

	// should not happen unless server's logic is changed
	return nil, errors.New("unexpected error - check server response codes")
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
}

func (c *Client) CreateUserLoginFile(username, password string) error {

	file, err := os.OpenFile(config.ClientCfg.UserData, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0600)
	if err != nil {
		return fmt.Errorf("error when trying to create/open userdata file:%w", err)
	}
	defer file.Close()

	hashedpassword, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return fmt.Errorf("error when trying to hash password:%w", err)
	}

	//writing to the file
	_, err = fmt.Fprintf(file, `{"login":"%s","password":"%s"}`, username, hashedpassword)
	if err != nil {
		return fmt.Errorf("error when trying to write into file:%w", err)
	}

	return nil
}
