package auth

import (
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

var InvalidSecretFileMsg = "please specify secret key file path"

// Returns new secret request with default client
func WithDefaultClient(teamId, clientId, keyId, secretKeyPath string) (*Request, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	return createRequest(teamId, clientId, keyId, secretKeyPath, client)
}

// Returns new secret request with given client
func WithCustomClient(client httpClient, teamId, clientId, keyId, secretKeyPath string) (*Request, error) {
	return createRequest(teamId, clientId, keyId, secretKeyPath, client)
}

func createRequest(teamId, clientId, keyId, secretKeyPath string, client httpClient) (*Request, error) {

	if secretKeyPath == "" {
		return nil, errors.New(InvalidSecretFileMsg)
	}

	secretContent, err := ioutil.ReadFile(secretKeyPath)

	if err != nil {
		return nil, err
	}

	return &Request{
		TeamID:       teamId,
		ClientID:     clientId,
		KeyID:        keyId,
		ClientSecret: secretContent,
		HttpClient:   client,
	}, nil
}
