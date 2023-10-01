package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"os"
)

func Encrypt(plainText string) (string, error) {
	secret := []byte(os.Getenv("AES_SECRET"))
	block, err := aes.NewCipher(secret)
	if err != nil {
		return "", err
	}
	print(len(plainText))

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	nonce := make([]byte, nonceSize)
	_, err = rand.Read(nonce)
	if err != nil {
		panic(err)
	}

	cipherText := gcm.Seal(nonce, nonce, []byte(plainText), nil)

	return hex.EncodeToString(cipherText), nil
}

func Decrypt(cipherText string) (string, error) {
	secret := []byte(os.Getenv("AES_SECRET"))

	cipherBytes, err := hex.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(secret)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	nonce, cipherBytes := cipherBytes[:nonceSize], cipherBytes[nonceSize:]

	plainText, err := gcm.Open(nil, nonce, cipherBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}
