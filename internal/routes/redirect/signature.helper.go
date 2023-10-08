package redirect

import (
	"errors"
	"shorty-urls-server/internal/internal/utils"
	"time"
)

type PasswordSignaturePayload struct {
	ID        string `json:"id"`
	ExpiresAt string `json:"expiresAt"`
	IssuedAt  string `json:"issuedAt"`
}

func generateSignature(id string) (string, error) {
	now := time.Now().UTC()
	payload := PasswordSignaturePayload{
		id,
		now.Add(time.Minute).String(),
		now.String(),
	}

	payloadString, err := utils.ConvertMapToString(payload)
	if err != nil {
		return "", err
	}

	encryptedString, err := utils.Encrypt(payloadString)
	if err != nil {
		return "", err
	}

	return encryptedString, nil
}

func validateSignature(signature string) (*PasswordSignaturePayload, error) {
	decryptedString, err := utils.Decrypt(signature)
	if err != nil {
		return nil, errors.New("Invalid signature")
	}

	payload, err := utils.ConvertStringToMap(decryptedString)
	if err != nil {
		return nil, errors.New("Invalid signature")
	}

	if payload["expiresAt"].(string) < time.Now().UTC().String() {
		return nil, errors.New("Signature expired")
	}

	payloadStruct := PasswordSignaturePayload{
		payload["id"].(string),
		payload["expiresAt"].(string),
		payload["issuedAt"].(string),
	}

	return &payloadStruct, nil
}
