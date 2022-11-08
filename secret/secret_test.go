package secret

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateSecret__emptyBlock(t *testing.T) {
	expected := "pem block is empty after decoding"

	req := Request{
		ClientID:     "com.example.app",
		TeamID:       "1234567890",
		KeyID:        "abc123def4",
		ClientSecret: "",
	}
	_, goterr := req.GenerateClientSecret()

	assert.Equal(t, expected, goterr.Error())
}
