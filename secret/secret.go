package secret

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	AUDIENCE = "https://appleid.apple.com"
)

type Request struct {
	// 10-char App Id prefix found in App identifiers section
	TeamID string

	//ClientID is the "Services ID" value that you get when navigating to your "sign in with Apple"-enabled service ID
	ClientID string

	// This is the ID of the private key
	KeyID string

	// This is the private key file (.p8). You can download it from apple portal
	ClientSecret []byte
}

// Returns new secret request
func New(teamId, clientId, keyId, secret string) *Request {
	return &Request{
		TeamID:       teamId,
		ClientID:     clientId,
		KeyID:        keyId,
		ClientSecret: []byte(secret),
	}
}

// GenerateClientSecret returns a secret used to validate server requests
// SecretRequest is required to generate secret. Method will throw error
// if data is empty or wrong.
func (req *Request) GenerateClientSecret() (string, error) {

	block, _ := pem.Decode(req.ClientSecret)
	if block == nil {
		return "", errors.New("pem block is empty after decoding")
	}

	prvKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, req.NewRegisteredClaims())
	token.Header["alg"] = "ES256"
	token.Header["kid"] = req.KeyID

	return token.SignedString(prvKey)
}

// NewRegisteredClaims generates jwt claims from SecretRequest.
func (req *Request) NewRegisteredClaims() *jwt.RegisteredClaims {

	// Current time
	now := time.Now()

	return &jwt.RegisteredClaims{
		Issuer:    req.TeamID,
		IssuedAt:  &jwt.NumericDate{Time: now},
		ExpiresAt: &jwt.NumericDate{Time: now.Add(time.Hour * 24 * 30)}, // 30 days
		Audience:  jwt.ClaimStrings{AUDIENCE},
		Subject:   req.ClientID,
	}
}
