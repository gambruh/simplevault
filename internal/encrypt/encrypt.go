package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

const filestorage = "./filestorage/filestorage"

func EncryptData(data, key []byte) ([]byte, error) {
	// Generate a new AES cipher block using the encryption key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create a new GCM (Galois/Counter Mode) cipher using the AES block
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Generate a random nonce (IV)
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Encrypt the plaintext using the GCM cipher
	ciphertext := gcm.Seal(nil, nonce, data, nil)

	// Concatenate the nonce and ciphertext and return the result
	encryptedData := append(nonce, ciphertext...)
	return encryptedData, nil
}

func StoreData(encryptedData []byte) error {
	// Encode the encrypted password in base64 for storage
	encodedData := base64.StdEncoding.EncodeToString(encryptedData)

	// Open the file in append mode or create the file if it doesn't exist
	file, err := os.OpenFile(filestorage, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the encoded password to the file
	_, err = fmt.Fprintln(file, encodedData)
	if err != nil {
		return err
	}

	return nil
}
