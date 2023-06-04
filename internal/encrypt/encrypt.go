package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
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
