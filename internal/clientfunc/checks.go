package clientfunc

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/gambruh/gophkeeper/internal/auth"
)

const userdata = "./userdata/user.json"

var (
	ErrLoginRequired    = errors.New("please login first")
	ErrNoCookieReturned = errors.New("server has not returned cookie")
	ErrWrongLoginData   = errors.New("wrong login data")
	ErrServerIsDown     = errors.New("server is down")
)

func GetUserData() (auth.LoginData, error) {

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
