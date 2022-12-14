// verify handles the receipt verification.
package receipt

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	// API endpoint for sandbox environment.
	SANDBOX_URL = "https://sandbox.itunes.apple.com/verifyReceipt"
	// API endpoint for production environment.
	PRODUCTION_URL = "https://buy.itunes.apple.com/verifyReceipt"
	// Request content-type for apple store.
	CONTENT_TYPE = "application/json; charset=utf-8"
)

// Returns new IAP request with defult client
func New() *Client {
	return &Client{
		HttpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Returns new IAP request with given client
func NewWithClient(client httpClient) *Client {
	return &Client{
		HttpClient: client,
	}
}

// Verify receipts and gets result from app store endpoints
func (client *Client) Verify(ctx context.Context, req IAPRequest) (response *IAPResponse, err error) {

	response, err = client.validateRequest(ctx, req, PRODUCTION_URL)
	if err != nil {
		return
	}

	// If receipt is from sandbox env
	if response.Status == 21007 {
		response, err = client.validateRequest(ctx, req, SANDBOX_URL)

		if err != nil {
			return
		}
	}

	return
}

// validates receipt request with production or sandbox urls
func (client *Client) validateRequest(ctx context.Context, req IAPRequest, url string) (*IAPResponse, error) {
	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(req); err != nil {
		return nil, err
	}

	newReq, err := http.NewRequestWithContext(ctx, "POST", url, b)
	if err != nil {
		return nil, err
	}

	newReq.Header.Add("content-type", CONTENT_TYPE)

	response, err := client.HttpClient.Do(newReq)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 500 {
		return nil, ErrAppStoreServer
	}

	return parseResponse(response)
}

// parse http response in IAP response
func parseResponse(resp *http.Response) (result *IAPResponse, err error) {

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	err = json.Unmarshal(buf, &result)
	if err != nil {
		return
	}

	return
}
