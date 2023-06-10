package helpers

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/gambruh/gophkeeper/internal/config"
	"github.com/gambruh/gophkeeper/internal/database"
	"github.com/gambruh/gophkeeper/internal/encrypt"
)

const privatekeyfile = "privatekey.pem"

var (
	ErrWrongFile = errors.New("file not found")
	ErrEmptyName = errors.New("metaname can't be empty")
)

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

func SplitFurther(input []string) (output []string) {

	if len(input) != 2 {
		return input
	}

	splitted := strings.Split(input[1], " ")
	output = input[:1]

	output = append(output, splitted...)

	return output
}

func EncryptNoteData(note database.Note, key []byte) (string, error) {
	// concatenating data to string
	noteStr := note.Name + "," + note.Text

	// encrypting the data
	encrypted, err := encrypt.EncryptData([]byte(noteStr), key)
	if err != nil {
		return "", err
	}
	// Encode the encrypted password in base64 for storage
	encodedData := base64.StdEncoding.EncodeToString(encrypted)

	return encodedData, nil
}

func DecryptNoteData(encrData database.EncryptedData, key []byte) (note database.Note, err error) {

	decodedData, err := base64.StdEncoding.DecodeString(encrData.Data)
	if err != nil {
		return database.Note{}, err
	}

	decryptedData, err := encrypt.DecryptData(decodedData, key)
	if err != nil {
		return database.Note{}, err
	}

	dst := string(decryptedData)

	noteArr := strings.Split(dst, ",")

	for i, data := range noteArr {
		switch i {
		case 0:
			note.Name = encrData.Name
		case 1:
			note.Text = data
		}
	}
	return note, nil
}

func ReadBinaryFile(filename string) ([]byte, error) {

	file, err := os.Open(filename)
	if err != nil {
		return nil, ErrWrongFile
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("error when reading file:%w", err)
	}

	return data, nil
}

func PrepareBinary(binaryname, sendfolder string) (newbinary database.Binary, err error) {
	if strings.HasPrefix(sendfolder, "/") {
		sendfolder = strings.TrimLeft(sendfolder, "/")
		sendfolder = "./" + sendfolder
	} else if strings.HasPrefix(sendfolder, "./") {
		// do nothing
	} else {
		sendfolder = "./" + sendfolder
	}

	if len(binaryname) == 0 {
		return database.Binary{}, ErrEmptyName
	}

	data, err := ReadBinaryFile(sendfolder + "/" + binaryname)
	if err != nil {
		return database.Binary{}, nil
	}

	newbinary.Name = binaryname
	newbinary.Data = data

	return newbinary, nil
}
