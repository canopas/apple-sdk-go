package auth

import "net/http"

// Returns new secret request
func NewClient(teamId, clientId, keyId, secret string) *Request {
	return &Request{
		TeamID:       teamId,
		ClientID:     clientId,
		KeyID:        keyId,
		ClientSecret: []byte(secret),
		HttpClient:   &http.Client{},
	}
}
