package encrypt

import (
	"bytes"
	"testing"
)

func TestEncryptDataDecryptData(t *testing.T) {
	key := []byte("0123456789abcdef")
	plaintext := []byte("Hello, World!")

	encryptedData, err := EncryptData(plaintext, key)
	if err != nil {
		t.Errorf("Error encrypting data: %v", err)
	}

	decryptedData, err := DecryptData(encryptedData, key)
	if err != nil {
		t.Errorf("Error decrypting data: %v", err)
	}

	if !bytes.Equal(plaintext, decryptedData) {
		t.Errorf("Decrypted data doesn't match the original plaintext")
	}
}
