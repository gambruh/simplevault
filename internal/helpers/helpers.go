package helpers

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net"
	"strings"

	"github.com/gambruh/gophkeeper/internal/config"
	"github.com/gambruh/gophkeeper/internal/database"
	"github.com/gambruh/gophkeeper/internal/encrypt"
)

const privatekeyfile = "privatekey.pem"

func CompareTwoMaps(mapServer, mapLocal map[string]struct{}) (toUpload map[string]struct{}, toDownload map[string]struct{}) {
	toUpload = make(map[string]struct{})
	toDownload = make(map[string]struct{})

	for iname := range mapServer {
		if _, ok := mapLocal[iname]; !ok {
			toDownload[iname] = struct{}{}
		}
	}

	for iname := range mapLocal {
		if _, ok := mapServer[iname]; !ok {
			toUpload[iname] = struct{}{}
		}
	}

	return toUpload, toDownload
}

func CreateMapFromList(list []string) (outputMap map[string]struct{}) {
	outputMap = make(map[string]struct{})

	for _, item := range list {
		outputMap[item] = struct{}{}
	}

	return outputMap
}

func CreateConfigTLS() (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(config.Cfg.Certificate, privatekeyfile)
	if err != nil {
		return nil, fmt.Errorf("error when loading tls certificate and key: %s", err)
	}

	tlsconfig := &tls.Config{
		GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
			fmt.Printf("Received TLS handshake request from client %s:%d\n", info.Conn.RemoteAddr().String(), info.Conn.RemoteAddr().(*net.TCPAddr).Port)
			cert, err := tls.LoadX509KeyPair(config.Cfg.Certificate, privatekeyfile)
			if err != nil {
				return nil, fmt.Errorf("error when loading KeyPair in CreateConfigTLS: %s", err)
			}
			return &cert, nil
		},
		Certificates: []tls.Certificate{cert},
	}
	return tlsconfig, nil
}

func EncryptCardData(card database.Card, key []byte) (string, error) {
	// concatenating card to string
	cardStr := card.Cardname + "," + card.Number + "," + card.Name + "," + card.Surname + "," + card.ValidTill + "," + card.Code

	// encrypting the card data
	encrypted, err := encrypt.EncryptData([]byte(cardStr), key)
	if err != nil {
		return "", err
	}
	// Encode the encrypted password in base64 for storage
	encodedData := base64.StdEncoding.EncodeToString(encrypted)

	return encodedData, nil
}

func DecryptCardData(encrCard database.EncryptedCard, key []byte) (database.Card, error) {
	var card database.Card
	decodedData, err := base64.StdEncoding.DecodeString(encrCard.Data)
	if err != nil {
		return database.Card{}, err
	}

	decryptedData, err := encrypt.DecryptData(decodedData, key)
	if err != nil {
		return database.Card{}, err
	}

	dst := string(decryptedData)

	cardArr := strings.Split(dst, ",")

	for i, data := range cardArr {
		switch i {
		case 0:
			card.Cardname = encrCard.Cardname
		case 1:
			card.Number = data
		case 2:
			card.Name = data
		case 3:
			card.Surname = data
		case 4:
			card.ValidTill = data
		case 5:
			card.Code = data
		}
	}
	return card, nil
}

func EncryptLoginCredsData(logincred database.LoginCreds, key []byte) (string, error) {
	// concatenating data to string
	logincredStr := logincred.Name + "," + logincred.Site + "," + logincred.Login + "," + logincred.Password

	// encrypting the data
	encrypted, err := encrypt.EncryptData([]byte(logincredStr), key)
	if err != nil {
		return "", err
	}
	// Encode the encrypted password in base64 for storage
	encodedData := base64.StdEncoding.EncodeToString(encrypted)

	return encodedData, nil
}

func DecryptLoginCredsData(encrData database.EncryptedData, key []byte) (database.LoginCreds, error) {
	var logincred database.LoginCreds

	decodedData, err := base64.StdEncoding.DecodeString(encrData.Data)
	if err != nil {
		return database.LoginCreds{}, err
	}

	decryptedData, err := encrypt.DecryptData(decodedData, key)
	if err != nil {
		return database.LoginCreds{}, err
	}

	dst := string(decryptedData)

	logincredArr := strings.Split(dst, ",")

	for i, data := range logincredArr {
		switch i {
		case 0:
			logincred.Name = encrData.Name
		case 1:
			logincred.Site = data
		case 2:
			logincred.Login = data
		case 3:
			logincred.Password = data
		}
	}
	return logincred, nil
}
