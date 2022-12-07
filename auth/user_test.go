package auth

import (
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockedClaims struct {
	mock.Mock
}

// Mocked function PostForm that does not call any server, just return the expected response.
func (m *MockedClaims) GetClaims(idToken string) (claims jwt.MapClaims, err error) {
	user := getUser()
	claims = jwt.MapClaims{
		"sub":              user.ID,
		"email":            user.Email,
		"email_verified":   user.EmailVerified,
		"is_private_email": user.IsPrivateEmail,
		"real_user_status": user.RealUserStatus,
	}
	return claims, nil
}

var resp TokenResponse

func TestGetUser(t *testing.T) {
	resp.Claims = new(MockedClaims)
	expected := getUser()
	got, err := resp.GetUser()

	assert.Equal(t, nil, err)
	assert.Equal(t, expected, got)
}

func TestUniqueID(t *testing.T) {
	resp.Claims = new(MockedClaims)
	expected := getUser().ID
	got, err := resp.UniqueID()

	assert.Equal(t, nil, err)
	assert.Equal(t, expected, got)

}

func TestEmail(t *testing.T) {
	resp.Claims = new(MockedClaims)
	expected := getUser().Email
	got, err := resp.Email()

	assert.Equal(t, nil, err)
	assert.Equal(t, expected, got)
}

func TestRealUserStatus(t *testing.T) {
	resp.Claims = new(MockedClaims)
	expected := getUser().RealUserStatus
	got, err := resp.RealUserStatus()

	assert.Equal(t, nil, err)
	assert.Equal(t, expected, got)

}

func getUser() *User {
	return &User{
		ID:             "123456",
		Email:          "john.doe@gmail.com",
		EmailVerified:  true,
		IsPrivateEmail: true,
		RealUserStatus: 2,
	}
}
