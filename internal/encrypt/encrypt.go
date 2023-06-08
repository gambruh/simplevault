package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
)

func EncryptData(data, key []byte) ([]byte, error) {

	aesblock, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		panic(err)
	}

	// создаём вектор инициализации
	nonce := key[len(key)-aesgcm.NonceSize():]

	// Encrypt the plaintext using the GCM cipher
	encryptedData := aesgcm.Seal(nil, nonce, data, nil)

	return encryptedData, nil
}

func DecryptData(encryptedData, key []byte) (decryptedData []byte, err error) {

	aesblock, err := aes.NewCipher(key[:])
	if err != nil {
		panic(err)
	}
	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		panic(err)
	}

	// создаём вектор инициализации
	nonce := key[len(key)-aesgcm.NonceSize():]

	//Decrypt the data
	decryptedData, err = aesgcm.Open(nil, nonce, encryptedData, nil) // расшифровываем
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	return decryptedData, nil
}

func DecryptFromString(s string, key []byte) (decryptedData []byte, err error) {
	dst, err := hex.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("can't decode string in DecryptFromString:%w", err)
	}
	return DecryptData(dst, key)
}

func EncryptFromString(s string, key []byte) (decryptedData []byte, err error) {
	dst, err := hex.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("can't decode string in EncryptFromString:%w", err)
	}
	return EncryptData(dst, key)
}
