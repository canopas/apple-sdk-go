package secret

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	expected := &Request{
		ClientID:     "com.example.app",
		TeamID:       "1234567890",
		KeyID:        "abc123def4",
		ClientSecret: []byte{},
	}

	got := New("1234567890", "com.example.app", "abc123def4", "")

	assert.Equal(t, expected, got)
}

func TestGenerateClientSecret__emptyBlock(t *testing.T) {
	expected := "pem block is empty after decoding"

	req := Request{
		ClientID:     "com.example.app",
		TeamID:       "1234567890",
		KeyID:        "abc123def4",
		ClientSecret: []byte{},
	}
	_, goterr := req.GenerateClientSecret()

	assert.Equal(t, expected, goterr.Error())
}
