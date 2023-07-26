// Package helpers provide misc functions used both by server and client
package helpers

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gambruh/gophkeeper/internal/encrypt"
	"github.com/gambruh/gophkeeper/internal/storage"
)

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

func EncryptCardData(card storage.Card, key []byte) (string, error) {
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

// DecryptCardData returns storage.Card struct out of encrypted data received from database
func DecryptCardData(encrCard storage.EncryptedData, key []byte) (storage.Card, error) {
	var card storage.Card
	decodedData, err := base64.StdEncoding.DecodeString(encrCard.Data)
	if err != nil {
		return storage.Card{}, err
	}

	decryptedData, err := encrypt.DecryptData(decodedData, key)
	if err != nil {
		return storage.Card{}, err
	}

	dst := string(decryptedData)

	cardArr := strings.Split(dst, ",")

	for i, data := range cardArr {
		switch i {
		case 0:
			card.Cardname = encrCard.Name
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

func EncryptLoginCredsData(logincred storage.LoginCreds, key []byte) (string, error) {
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

// DecryptLoginCredsData returns storage.LoginCreds struct out of encrypted data received from database
func DecryptLoginCredsData(encrData storage.EncryptedData, key []byte) (storage.LoginCreds, error) {
	var logincred storage.LoginCreds

	decodedData, err := base64.StdEncoding.DecodeString(encrData.Data)
	if err != nil {
		return storage.LoginCreds{}, err
	}

	decryptedData, err := encrypt.DecryptData(decodedData, key)
	if err != nil {
		return storage.LoginCreds{}, err
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

// SplitFurther is a helper function to work with commands in CLI
func SplitFurther(input []string) (output []string) {

	if len(input) != 2 {
		return input
	}

	splitted := strings.Split(input[1], " ")
	output = input[:1]

	output = append(output, splitted...)

	return output
}

// EncryptNoteData encrypts storage.Note and saves it to a hexadecimal string to be stored in a database
func EncryptNoteData(note storage.Note, key []byte) (string, error) {
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

func DecryptNoteData(encrData storage.EncryptedData, key []byte) (note storage.Note, err error) {

	decodedData, err := base64.StdEncoding.DecodeString(encrData.Data)
	if err != nil {
		return storage.Note{}, err
	}

	decryptedData, err := encrypt.DecryptData(decodedData, key)
	if err != nil {
		return storage.Note{}, err
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

// ReadBinaryFile reads data from binary file and returns its contents
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

// PrepareBinary returns storage.Binary struct out of binary file which can be saved to a storage
func PrepareBinary(binaryname, sendfolder string) (newbinary storage.Binary, err error) {
	if strings.HasPrefix(sendfolder, "/") {
		sendfolder = strings.TrimLeft(sendfolder, "/")
		sendfolder = "./" + sendfolder
	} else if strings.HasPrefix(sendfolder, "./") {
		// do nothing
	} else {
		sendfolder = "./" + sendfolder
	}

	if len(binaryname) == 0 {
		return storage.Binary{}, ErrEmptyName
	}

	data, err := ReadBinaryFile(sendfolder + "/" + binaryname)
	if err != nil {
		return storage.Binary{}, nil
	}

	newbinary.Name = binaryname
	newbinary.Data = data

	return newbinary, nil
}
