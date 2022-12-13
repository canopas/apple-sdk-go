package auth

import (
	"io/ioutil"
	"net/http"
)

// Returns new secret request
func NewClient(teamId, clientId, keyId, secretKeyPath string) (*Request, error) {
	secretContent, err := ioutil.ReadFile(secretKeyPath)

	if err != nil {
		return nil, err
	}

	return &Request{
		TeamID:       teamId,
		ClientID:     clientId,
		KeyID:        keyId,
		ClientSecret: secretContent,
		HttpClient:   &http.Client{},
	}, nil
}
