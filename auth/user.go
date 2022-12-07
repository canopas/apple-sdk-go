// user retrive the information about validated apple user from idToken
package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt"
)

type claims interface {
	GetClaims(idToken string) (jwt.MapClaims, error)
}

type Claims struct{}

// User will have the information of authenticated user of Apple.
type User struct {
	// The unique identifier for the user (sub).
	ID string `json:"id"`

	// A string value that represents the user’s email address.
	// The email address is either the user’s real email address or the proxy address,
	// depending on their private email relay service.
	Email string `json:"email"`

	// A string or Boolean value that indicates whether the service verifies the email.
	EmailVerified bool `json:"email_verified"`

	// A string or Boolean value that indicates whether the email
	// that the user shares is the proxy address.
	// The value can either be a string ("true" or "false") or a Boolean (true or false).
	IsPrivateEmail bool `json:"is_private_email"`

	// An Integer value that indicates whether the user appears to be a real person.
	// Use the value of this claim to mitigate fraud.
	// The possible values are: 0 (or Unsupported), 1 (or Unknown), 2 (or LikelyReal).
	RealUserStatus int `json:"real_user_status"`
}

// GetClaims decodes the idToken and returns the claims
func (c *Claims) GetClaims(idToken string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(idToken, claims, func(token *jwt.Token) (interface{}, error) {
		return "", nil
	})
	if err != nil {
		return nil, err
	}

	return claims, nil
}

// UniqueID returns the unique subject ID to identify the user
func (resp *TokenResponse) UniqueID() (string, error) {
	claims, err := resp.Claims.GetClaims(resp.IDToken)
	return fmt.Sprintf("%v", claims["sub"]), err
}

// Email returns the user email
func (resp *TokenResponse) Email() (string, error) {
	claims, err := resp.Claims.GetClaims(resp.IDToken)
	return fmt.Sprintf("%v", claims["email"]), err
}

// RealUserStatus returns whether the user appears to be a real person.
// The possible values are: 0 (or Unsupported), 1 (or Unknown), 2 (or LikelyReal).
func (resp *TokenResponse) RealUserStatus() (int, error) {
	claims, err := resp.Claims.GetClaims(resp.IDToken)
	return claims["real_user_status"].(int), err
}

// GetUser will get claims, and returns the user using claims
func (resp *TokenResponse) GetUser() (*User, error) {
	claims, err := resp.Claims.GetClaims(resp.IDToken)
	if err != nil {
		return nil, err
	}

	var user User

	if sub, ok := claims["sub"].(string); ok {
		user.ID = sub
	}

	if email, ok := claims["email"].(string); ok {
		user.Email = email
	}

	if realUserStatus, ok := claims["real_user_status"].(int); ok {
		user.RealUserStatus = realUserStatus
	}

	if emailVerified, ok := claims["email_verified"].(bool); ok {
		user.EmailVerified = emailVerified
	}

	if isPrivateEmail, ok := claims["is_private_email"].(bool); ok {
		user.IsPrivateEmail = isPrivateEmail
	}

	return &user, nil
}
