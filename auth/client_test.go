package auth

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateRequest(t *testing.T) {
	gotResp, gotErr := createRequest("1234567890", "com.example.app", "abc123def4", "", &http.Client{})
	expectedErr := errors.New(InvalidSecretFileMsg)

	assert.Equal(t, expectedErr, gotErr)
	assert.NotEqual(t, nil, gotResp)
}
