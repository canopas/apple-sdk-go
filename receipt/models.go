// models have receipt JSON information in structs.
package receipt

import (
	"encoding/json"
	"net/http"
)

type httpClient interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

type (
	numericString string

	// Receipt client with http client
	Client struct {
		HttpClient httpClient
	}

	// The JSON contents you submit with the request to the App Store.
	// https://developer.apple.com/documentation/appstorereceipts/requestbody
	IAPRequest struct {
		// (Required) The Base64-encoded receipt data.
		ReceiptData string `json:"receipt-data"`

		// Your app’s shared secret, which is a hexadecimal string.
		Password string `json:"password,omitempty"`

		// Only used for iOS7 style app receipts that contain auto-renewable or non-renewing subscriptions.
		// Set this value to true for the response to include only the latest renewal transaction for any subscriptions.
		// Use this field only for app receipts that contain auto-renewable subscriptions.
		ExcludeOldTransactions bool `json:"exclude-old-transactions"`
	}

	// The JSON data that returns in the response from the App Store.
	// https://developer.apple.com/documentation/appstorereceipts/responsebody
	IAPResponse struct {
		// Either 0 if the receipt is valid, or a status code if there’s an error.
		// The status code reflects the status of the app receipt as a whole
		Status int `json:"status"`

		// The environment the system generates the receipt for.
		// Possible values: Sandbox, Production
		Environment string `json:"environment"`

		// A JSON representation of the receipt that you send for verification.
		Receipt Receipt `json:"receipt"`

		// An array that contains all in-app purchase transactions.
		// This excludes transactions for consumable products that your app marks as finished.
		// This only returns for receipts that contain auto-renewable subscriptions.
		LatestReceiptInfo []InApp `json:"latest_receipt_info,omitempty"`

		// The latest Base64-encoded app receipt.
		// This only returns for receipts that contain auto-renewable subscriptions.
		LatestReceipt string `json:"latest_receipt,omitempty"`

		// An array where each element contains the pending renewal information for each auto-renewable subscription the product_id identifies.
		// This only returns for app receipts that contain auto-renewable subscriptions.
		PendingRenewalInfo []PendingRenewalInfo `json:"pending_renewal_info,omitempty"`

		IsRetryable bool `json:"is-retryable,omitempty"`
	}

	// An array that contains the in-app purchase receipt fields for all in-app purchase transactions.
	// https://developer.apple.com/documentation/appstorereceipts/responsebody/receipt/in_app
	InApp struct {
		// The number of consumable products purchased.
		// This value corresponds to the quantity property of the SKPayment object stored in the transaction's payment property.
		// The value is usually "1" unless modified with a mutable payment. The maximum value is 10.
		Quantity string `json:"quantity"`

		// The unique identifier of the product purchased.
		ProductID string `json:"product_id"`

		// A unique identifier for a transaction such as a purchase, restore, or renewal.
		TransactionID string `json:"transaction_id"`

		// The transaction identifier of the original purchase.
		OriginalTransactionID string `json:"original_transaction_id"`

		// A unique identifier for purchase events across devices, including subscription-renewal events.
		// This value is the primary key for identifying subscription purchases.
		WebOrderLineItemID string `json:"web_order_line_item_id,omitempty"`

		// The identifier of the promotional offer for an auto-renewable subscription that the user redeems.
		PromotionalOfferID string `json:"promotional_offer_id"`

		// An indication of whether a subscription is in the free trial period.
		// https://developer.apple.com/documentation/appstorereceipts/is_trial_period
		IsTrialPeriod string `json:"is_trial_period"`

		// An indicator of whether an auto-renewable subscription is in the introductory price period.
		// https://developer.apple.com/documentation/appstorereceipts/is_in_intro_offer_period.
		IsInIntroOfferPeriod string `json:"is_in_intro_offer_period,omitempty"`

		ExpiresDate
		PurchaseDate
		OriginalPurchaseDate
		CancellationDate

		// The reason for a refunded or revoked transaction.
		// A value of "1" indicates that the customer canceled their transaction due to an actual or perceived issue within your app.
		// A value of "0" indicates that the transaction was canceled for another reason;
		// for example, if the customer made the purchase accidentally.
		// Possible values: 1, 0
		CancellationReason string `json:"cancellation_reason,omitempty"`
	}

	// The decoded version of the encoded receipt data that you send with the request to the App Store.
	// https://developer.apple.com/documentation/appstorereceipts/responsebody/receipt
	Receipt struct {
		// The type of receipt generated.
		// The value corresponds to the environment in which the app or VPP purchase was made.
		// Possible values: Production, ProductionVPP, ProductionSandbox, ProductionVPPSandbox
		ReceiptType string `json:"receipt_type"`

		// See app_item_id.
		AdamID int64 `json:"adam_id"`

		// Generated by App Store Connect and used by the App Store to uniquely identify the app purchased.
		// Apps are assigned this identifier only in production. Treat this value as a 64-bit long integer.
		AppItemID numericString `json:"app_item_id"`

		// The bundle identifier for the app to which the receipt belongs.
		// You provide this string on App Store Connect.
		// This corresponds to the value of CFBundleIdentifier in the Info.plist file of the app.
		BundleID string `json:"bundle_id"`

		// The app’s version number. The app's version number corresponds to the value of CFBundleVersion (in iOS) or CFBundleShortVersionString (in macOS) in the Info.plist.
		// In production, this value is the current version of the app on the device based on the receipt_creation_date_ms.
		// In the sandbox, the value is always "1.0".
		ApplicationVersion string `json:"application_version"`

		// A unique identifier for the app download transaction.
		DownloadID int64 `json:"download_id"`

		// An arbitrary number that identifies a revision of your app. In the sandbox, this key's value is "0".
		VersionExternalIdentifier numericString `json:"version_external_identifier"`

		// The version of the app that the user originally purchased.
		// This value does not change, and corresponds to the value of CFBundleVersion (in iOS) or CFBundleShortVersionString (in macOS) in the Info.plist file of the original purchase.
		// In the sandbox environment, the value is always "1.0".
		OriginalApplicationVersion string `json:"original_application_version"`

		// An array that contains the in-app purchase receipt fields for all in-app purchase transactions.
		InApp []InApp `json:"in_app"`

		ReceiptCreationDate
		RequestDate
		OriginalPurchaseDate
		PreorderDate
		ExpiresDate
	}

	// An array of elements that refers to open or failed auto-renewable subscription renewals.
	// https://developer.apple.com/documentation/appstorereceipts/responsebody/pending_renewal_info
	PendingRenewalInfo struct {
		// The reason a subscription expired.
		// This field is present only for a receipt that contains an expired auto-renewable subscription.
		SubscriptionExpirationIntent string `json:"expiration_intent"`

		// The value for this key corresponds to the productIdentifier property of
		// the product that the customer’s subscription renews.
		SubscriptionAutoRenewProductID string `json:"auto_renew_product_id"`

		// A flag that indicates Apple is attempting to renew an expired subscription automatically.
		// This field is present only if an auto-renewable subscription is in the billing retry state.
		// https://developer.apple.com/documentation/appstorereceipts/is_in_billing_retry_period
		SubscriptionRetryFlag string `json:"is_in_billing_retry_period"`

		// The current renewal status for the auto-renewable subscription.
		SubscriptionAutoRenewStatus string `json:"auto_renew_status"`

		// The price consent status for an auto-renewable subscription price increase that requires customer consent.
		// This field is present only if the App Store requested customer consent for a price increase that requires customer consent.
		// The default value is "0" and changes to "1" if the customer consents.
		// Possible values: 1, 0
		SubscriptionPriceConsentStatus string `json:"price_consent_status"`

		// The unique identifier of the product purchased.
		ProductID string `json:"product_id"`

		// The transaction identifier of the original purchase.
		OriginalTransactionID string `json:"original_transaction_id"`

		// The reference name of a subscription offer that you configured in App Store Connect.
		// This field is present when a customer redeemed a subscription offer code.
		// https://developer.apple.com/documentation/appstorereceipts/offer_code_ref_name
		OfferCodeRefName string `json:"offer_code_ref_name,omitempty"`

		// The identifier of the promotional offer for an auto-renewable subscription that the user redeemed.
		// You provide this value in the Promotional Offer Identifier field when you create the promotional offer in App Store Connect.
		PromotionalOfferID string `json:"promotional_offer_id,omitempty"`

		// The status that indicates if an auto-renewable subscription is subject to a price increase.
		// The price increase status is "0" when the App Store has requested consent
		// for an auto-renewable subscription price increase that requires customer consent,
		// and the customer hasn't yet consented.

		// The price increase status is "1" if the customer has consented to a price increase
		// that requires customer consent.

		// The price increase status is also "1" if the App Store has notified the customer of
		// the auto-renewable subscription price increase that doesn't require customer consent.

		// Possible values: 1, 0
		PriceIncreaseStatus string `json:"price_increase_status,omitempty"`

		GracePeriodDate
	}

	// The date when the app receipt was created.
	ReceiptCreationDate struct {
		// The time the App Store generated the receipt, in a date-time format similar to ISO 8601.
		CreationDate string `json:"receipt_creation_date"`

		// The time the App Store generated the receipt, in UNIX epoch time format, in milliseconds.
		// Use this time format for processing dates. This value does not change.
		CreationDateMS string `json:"receipt_creation_date_ms"`

		// The time the App Store generated the receipt, in the Pacific Time zone.
		CreationDatePST string `json:"receipt_creation_date_pst"`
	}

	// The date and time of when the request was sent
	RequestDate struct {
		// The time the request to the verifyReceipt endpoint was processed and the response was generated,
		// in a date-time format similar to ISO 8601.
		RequestDate string `json:"request_date"`

		// The time the request to the verifyReceipt endpoint was processed and the response was generated, in UNIX epoch time format, in milliseconds.
		// Use this time format for processing dates.
		RequestDateMS string `json:"request_date_ms"`

		// The time the request to the verifyReceipt endpoint was processed and the response was generated,
		// in the Pacific Time zone.
		RequestDatePST string `json:"request_date_pst"`
	}

	// The date and time of when the item was purchased
	PurchaseDate struct {
		// The time the App Store charged the user's account for a purchased or restored product,
		// or the time the App Store charged the user’s account for a subscription purchase or renewal
		// after a lapse, in a date-time format similar to ISO 8601.
		PurchaseDate string `json:"purchase_date"`

		// For consumable, non-consumable, and non-renewing subscription products,
		// the time the App Store charged the user's account for a purchased or restored product,
		// in the UNIX epoch time format, in milliseconds.
		// For auto-renewable subscriptions, the time the App Store charged the user’s account
		// for a subscription purchase or renewal after a lapse, in the UNIX epoch time format, in milliseconds.
		// Use this time format for processing dates.
		PurchaseDateMS string `json:"purchase_date_ms"`

		// The time the App Store charged the user's account for a purchased or restored product,
		// or the time the App Store charged the user’s account
		// for a subscription purchase or renewal after a lapse, in the Pacific Time zone.
		PurchaseDatePST string `json:"purchase_date_pst"`
	}

	// The beginning of the subscription period
	OriginalPurchaseDate struct {
		// The time of the original app purchase, in a date-time format similar to ISO 8601.
		OriginalPurchaseDate string `json:"original_purchase_date"`

		// The time of the original app purchase, in UNIX epoch time format, in milliseconds.
		// Use this time format for processing dates.
		OriginalPurchaseDateMS string `json:"original_purchase_date_ms"`

		// The time of the original app purchase, in the Pacific Time zone.
		OriginalPurchaseDatePST string `json:"original_purchase_date_pst"`
	}

	// The date and time of the pre-order
	PreorderDate struct {
		// The time the user ordered the app available for pre-order,
		// in a date-time format similar to ISO 8601.
		PreorderDate string `json:"preorder_date"`

		// The time the user ordered the app available for pre-order, in UNIX epoch time format, in milliseconds.
		// This field is only present if the user pre-orders the app.
		// Use this time format for processing dates.
		PreorderDateMS string `json:"preorder_date_ms"`

		// The time the user ordered the app available for pre-order, in the Pacific Time zone.
		PreorderDatePST string `json:"preorder_date_pst"`
	}

	// The expiration date for the subscription
	ExpiresDate struct {
		// The time the receipt expires for apps purchased through the Volume Purchase Program,
		// in a date-time format similar to the ISO 8601.
		ExpiresDate string `json:"expires_date,omitempty"`

		// The time the receipt expires for apps purchased through the Volume Purchase Program,
		// in UNIX epoch time format, in milliseconds.
		// If this key is not present for apps purchased through the Volume Purchase Program,
		// the receipt does not expire.
		// Use this time format for processing dates.
		ExpiresDateMS string `json:"expires_date_ms,omitempty"`

		// The time the receipt expires for apps purchased through the Volume Purchase Program,
		// in the Pacific Time zone.
		ExpiresDatePST string `json:"expires_date_pst,omitempty"`
	}

	// The time and date of the cancellation by Apple customer support
	CancellationDate struct {
		// The time the App Store refunded a transaction or revoked it from family sharing,
		// in a date-time format similar to the ISO 8601.
		// This field is present only for refunded or revoked transactions.
		CancellationDate string `json:"cancellation_date,omitempty"`

		// The time the App Store refunded a transaction or revoked it from family sharing,
		// in UNIX epoch time format, in milliseconds.
		// This field is present only for refunded or revoked transactions.
		// Use this time format for processing dates. See cancellation_date_ms for more information.
		CancellationDateMS string `json:"cancellation_date_ms,omitempty"`

		// The time the App Store refunded a transaction or revoked it from family sharing, in the Pacific Time zone.
		// This field is present only for refunded or revoked transactions.
		CancellationDatePST string `json:"cancellation_date_pst,omitempty"`
	}

	// Grace period date for the subscription
	GracePeriodDate struct {
		// The time at which the grace period for subscription renewals expires,
		// in a date-time format similar to the ISO 8601.
		GracePeriodDate string `json:"grace_period_expires_date,omitempty"`

		// The time at which the grace period for subscription renewals expires, in UNIX epoch time format, in milliseconds.
		// This key is present only for apps that have Billing Grace Period enabled and when the user experiences a billing error at the time of renewal.
		// Use this time format for processing dates.
		GracePeriodDateMS string `json:"grace_period_expires_date_ms,omitempty"`

		// The time at which the grace period for subscription renewals expires, in the Pacific Time zone.
		GracePeriodDatePST string `json:"grace_period_expires_date_pst,omitempty"`
	}
)

func (n *numericString) UnmarshalJSON(b []byte) error {
	var number json.Number
	if err := json.Unmarshal(b, &number); err != nil {
		return err
	}
	*n = numericString(number.String())
	return nil
}
