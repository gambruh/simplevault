package clientfunc

import "errors"

var (
	ErrLoginRequired    = errors.New("please login first")
	ErrNoCookieReturned = errors.New("server has not returned cookie")
	ErrWrongLoginData   = errors.New("wrong login data")
	ErrServerIsDown     = errors.New("server is down")
	ErrUsernameIsTaken  = errors.New("username is taken")
	ErrMetanameIsTaken  = errors.New("metaname(cardname) already in use, provide new one")
	ErrDataNotFound     = errors.New("data not found")
	ErrBadRequest       = errors.New("bad request")
)
