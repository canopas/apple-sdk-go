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
	// Endpoint for sandbox environment.
	SANDBOX_URL = "https://sandbox.itunes.apple.com/verifyReceipt"
	// Endpoint for production environment.
	PRODUCTION_URL = "https://buy.itunes.apple.com/verifyReceipt"
	// Request content-type for apple store.
	CONTENT_TYPE = "application/json; charset=utf-8"
)

type (
	numericString string

	Client struct {
		HttpClient httpClient
	}

	IAPRequest struct {
		ReceiptData string `json:"receipt-data"`
		Password    string `json:"password,omitempty"`
	}

	IAPResponse struct {
		Status             int                  `json:"status"`
		Environment        string               `json:"environment"`
		Receipt            Receipt              `json:"receipt"`
		LatestReceiptInfo  []InApp              `json:"latest_receipt_info,omitempty"`
		LatestReceipt      string               `json:"latest_receipt,omitempty"`
		PendingRenewalInfo []PendingRenewalInfo `json:"pending_renewal_info,omitempty"`
		IsRetryable        bool                 `json:"is-retryable,omitempty"`
	}

	// The InApp type has the receipt attributes
	InApp struct {
		Quantity                    string `json:"quantity"`
		ProductID                   string `json:"product_id"`
		TransactionID               string `json:"transaction_id"`
		OriginalTransactionID       string `json:"original_transaction_id"`
		PromotionalOfferID          string `json:"promotional_offer_id"`
		SubscriptionGroupIdentifier string `json:"subscription_group_identifier"`
		OfferCodeRefName            string `json:"offer_code_ref_name,omitempty"`
		ExpiresDateMS               string `json:"expires_date_ms,omitempty"`
		PurchaseDateMS              string `json:"purchase_date_ms"`
		CancellationDateMS          string `json:"cancellation_date_ms,omitempty"`
		CancellationReason          string `json:"cancellation_reason,omitempty"`
	}

	// Receipt data
	Receipt struct {
		ReceiptType                string        `json:"receipt_type"`
		AdamID                     int64         `json:"adam_id"`
		AppItemID                  numericString `json:"app_item_id"`
		BundleID                   string        `json:"bundle_id"`
		ApplicationVersion         string        `json:"application_version"`
		DownloadID                 int64         `json:"download_id"`
		OriginalApplicationVersion string        `json:"original_application_version"`
		InApp                      []InApp       `json:"in_app"`
	}

	PendingRenewalInfo struct {
		SubscriptionExpirationIntent   string `json:"expiration_intent"`
		SubscriptionAutoRenewProductID string `json:"auto_renew_product_id"`
		SubscriptionRetryFlag          string `json:"is_in_billing_retry_period"`
		SubscriptionAutoRenewStatus    string `json:"auto_renew_status"`
		SubscriptionPriceConsentStatus string `json:"price_consent_status"`
		ProductID                      string `json:"product_id"`
		OriginalTransactionID          string `json:"original_transaction_id"`
		OfferCodeRefName               string `json:"offer_code_ref_name,omitempty"`
		GracePeriodDateMS              string `json:"grace_period_expires_date_ms,omitempty"`
	}
)

type httpClient interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

func (n *numericString) UnmarshalJSON(b []byte) error {
	var number json.Number
	if err := json.Unmarshal(b, &number); err != nil {
		return err
	}
	*n = numericString(number.String())
	return nil
}

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

// Verify receipts and gets result
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
