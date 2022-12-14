package receipt

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
)

// Implementation of a mocked HTTP client.
type MockedHTTPClient struct {
	mock.Mock
}

// Mocked function PostForm that does not call any server, just return the expected response.
func (m *MockedHTTPClient) Do(req *http.Request) (resp *http.Response, err error) {
	return resp, ErrInvalidReceiptData
}

func TestValidateRequest(t *testing.T) {
	cli := Client{
		HttpClient: new(MockedHTTPClient),
	}
	gotResp, gotErr := cli.validateRequest(context.Background(), IAPRequest{ReceiptData: "", Password: ""}, SANDBOX_URL)

	assert.Equal(t, ErrInvalidReceiptData, gotErr)
	assert.NotEqual(t, nil, gotResp)
}

func TestVerify(t *testing.T) {
	cli := Client{
		HttpClient: new(MockedHTTPClient),
	}
	gotResp, gotErr := cli.Verify(context.Background(), IAPRequest{ReceiptData: "", Password: ""})

	assert.Equal(t, ErrInvalidReceiptData, gotErr)
	assert.NotEqual(t, nil, gotResp)
}

func TestHandleErrors(t *testing.T) {
	gotResp := HandleErrors(21000)
	assert.Equal(t, ErrInvalidJSON, gotResp)

	gotResp = HandleErrors(21004)
	assert.Equal(t, ErrInvalidSharedSecret, gotResp)

	gotResp = HandleErrors(21010)
	assert.Equal(t, ErrReceiptUnauthorized, gotResp)
}
