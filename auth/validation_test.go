package auth

import (
	"context"
	"errors"

	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Implementation of a mocked HTTP client.
type MockedHTTPClient struct {
	mock.Mock
}

// Mocked function PostForm that does not call any server, just return the expected response.
func (m *MockedHTTPClient) Do(req *http.Request) (resp *http.Response, err error) {
	return resp, errors.New(InvalidClientMsg)
}

func TestNew(t *testing.T) {
	expected := request()
	got := New("1234567890", "com.example.app", "abc123def4", "")
	assert.Equal(t, expected, got)
}

func TestErrorResponse(t *testing.T) {
	expected := errors.New("The requested scope is invalid.")
	got := errorResponse(ErrorResponse{
		Error: "invalid_scope",
	})

	assert.Equal(t, expected, got)
}

func TestGenerateClientSecret__emptyBlock(t *testing.T) {
	expected := "pem block is empty after decoding"

	_, goterr := request().GenerateClientSecret()

	assert.Equal(t, expected, goterr.Error())
}

func TestDoRequest(t *testing.T) {
	expectedErr := errors.New(InvalidClientMsg)

	vReq := request()
	vReq.HttpClient = new(MockedHTTPClient)
	gotResp, gotErr := vReq.doRequest(context.Background(), formData())

	assert.Equal(t, expectedErr, gotErr)
	assert.NotEqual(t, nil, gotResp)
}

func formData() url.Values {
	return url.Values{
		"client_id":     []string{"com.example.app"},
		"client_secret": []string{""},
		"grant_type":    []string{authGrantType},
		"code":          []string{"123456"},
	}
}

func request() *Request {
	return &Request{
		ClientID:     "com.example.app",
		TeamID:       "1234567890",
		KeyID:        "abc123def4",
		ClientSecret: []byte{},
		HttpClient:   &http.Client{},
	}
}
