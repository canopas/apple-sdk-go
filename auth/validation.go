// validation handles the sign in token validations.
package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

const (
	VALIDATION_URL = "https://appleid.apple.com/auth/token"
)

var (
	// Grant types for auth and refreshToken
	authGrantType         string = "authorization_code"
	refreshTokenGrantType string = "refresh_token"

	// The request is malformed, typically because it’s missing a parameter,
	// contains an unsupported parameter, includes multiple credentials,
	// or uses more than one mechanism for authenticating the client.
	InvalidRequest    string = "invalid_request"
	InvalidRequestMsg string = "The request is malformed, typically because it is missing a parameter, contains an unsupported parameter, includes multiple credentials, or uses more than one mechanism for authenticating the client."

	// The client authentication failed, typically due to a mismatched or invalid client identifier,
	// invalid client secret (expired token, malformed claims, or invalid signature), or mismatched or invalid redirect URI.
	InvalidClient    string = "invalid_client"
	InvalidClientMsg string = "The client authentication failed, typically due to a mismatched or invalid client identifier, invalid client secret (expired token, malformed claims, or invalid signature), or mismatched or invalid redirect URI."

	// The authorization grant or refresh token is invalid,
	// typically due to a mismatched or invalid client identifier,
	// invalid code (expired or previously used authorization code),
	// or invalid refresh token.
	InvalidGrant    string = "invalid_grant"
	InvalidGrantMsg string = "The authorization grant or refresh token is invalid, typically due to a mismatched or invalid client identifier, invalid code (expired or previously used authorization code), or invalid refresh token."

	// The client isn’t authorized to use this authorization grant type.
	UnauthorizedClient    string = "unauthorized_client"
	UnauthorizedClientMsg string = "The client is not authorized to use this authorization grant type."

	// The authenticated client isn’t authorized to use this grant type.
	UnsupportedGrantType    string = "unsupported_grant_type"
	UnsupportedGrantTypeMsg string = "The authenticated client is not authorized to use this grant type."

	// The requested scope is invalid.
	InvalidScope    string = "invalid_scope"
	InvalidScopeMsg string = "The requested scope is invalid."
)

type Validation interface {

	// Validates request using the authorization code received in an authorization
	// response sent to your app.
	// Returns accessToken, refreshToken, idToken
	ValidateCode(code string) (*TokenResponse, error)

	// Validate request using destinatio URI provided in authorization request
	// Returns accessToken, refreshToken, idToken
	ValidateCodeWithRedirectURI(code, redirectURI string) (*TokenResponse, error)

	// Validates given refresh token
	// Returns accessToken and idToken
	ValidateRefreshToken(refreshToken string) (*TokenResponse, error)
}

// Response after validation process from apple
type TokenResponse struct {

	// The refresh token used to regenerate new access tokens when validating an authorization code.
	// Store this token securely on your server.
	// The refresh token isn’t returned when validating an existing refresh token.
	RefreshToken string `json:"refresh_token"`

	// A token used to access allowed data,
	// such as generating and exchanging transfer identifiers during user migration
	AccessToken string `json:"access_token"`

	// The amount of time, in seconds, before the access token expires.
	ExpiresIn int `json:"expires_in"`

	// A JSON Web Token (JWT) that contains the user’s identity information.
	IDToken string `json:"id_token"`

	// The type of access token, which is always bearer.
	TokenType string `json:"token_type"`

	Claims claims
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// Validates request using the authorization code received in an authorization
// response sent to your app.
// Returns TokenResponse and error
func (req *Request) ValidateCode(code string) (*TokenResponse, error) {
	formData, err := req.newFormData(code, authGrantType, "", "")
	if err != nil {
		return nil, err
	}
	return req.doRequest(formData)
}

// Validate request using destinatio URI provided in authorization request
// Returns TokenResponse and error
func (req *Request) ValidateCodeWithRedirectURI(code, redirectURI string) (*TokenResponse, error) {
	formData, err := req.newFormData(code, authGrantType, redirectURI, "")
	if err != nil {
		return nil, err
	}
	return req.doRequest(formData)
}

// Validates given refresh token
// Returns TokenResponse and error
func (req *Request) ValidateRefreshToken(refreshToken string) (*TokenResponse, error) {
	formData, err := req.newFormData("", refreshTokenGrantType, "", refreshToken)
	if err != nil {
		return nil, err
	}
	return req.doRequest(formData)
}

func (req *Request) doRequest(formData url.Values) (*TokenResponse, error) {

	response, err := req.HttpClient.PostForm(VALIDATION_URL, formData)

	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.NewDecoder(response.Body).Decode(&errResp); err != nil {
			return nil, err
		}
		return nil, errorResponse(errResp)
	}

	var tokenResponse TokenResponse
	if err := json.NewDecoder(response.Body).Decode(&tokenResponse); err != nil {
		return nil, err
	}

	tokenResponse.Claims = &Claims{}

	return &tokenResponse, nil
}

// Prepare form data from given data
func (req *Request) newFormData(code, grantType, redirectURI, refreshToken string) (url.Values, error) {

	formData := make(url.Values)

	secret, err := req.GenerateClientSecret()

	if err != nil {
		return nil, err
	}

	formData.Add("client_id", req.ClientID)
	formData.Add("client_secret", secret)
	formData.Add("grant_type", grantType)

	if code != "" {
		formData.Add("code", code)
	}

	if redirectURI != "" {
		formData.Add("redirect_uri", redirectURI)
	}

	if refreshToken != "" {
		formData.Add("refresh_token", refreshToken)
	}

	return formData, nil
}

func errorResponse(err ErrorResponse) error {

	switch err.Error {
	case InvalidRequest:
		return errors.New(InvalidRequestMsg)

	case InvalidClient:
		return errors.New(InvalidClientMsg)

	case InvalidGrant:
		return errors.New(InvalidGrantMsg)

	case UnauthorizedClient:
		return errors.New(UnauthorizedClientMsg)

	case UnsupportedGrantType:
		return errors.New(UnsupportedGrantTypeMsg)

	case InvalidScope:
		return errors.New(InvalidScopeMsg)

	default:
		return errors.New("Unrecognized error: " + err.Error)
	}

}
