package receipt

import "errors"

// list of errors
var (
	ErrAppStoreServer          = errors.New("appStore server error")
	ErrInvalidJSON             = errors.New("the App Store could not read the JSON object you provided")
	ErrInvalidReceiptData      = errors.New("the data in the receipt-data property was malformed or missing")
	ErrReceiptUnauthenticated  = errors.New("the receipt could not be authenticated")
	ErrInvalidSharedSecret     = errors.New("the shared secret you provided does not match the shared secret on file for your account")
	ErrServerUnavailable       = errors.New("the receipt server is not currently available")
	ErrReceiptIsForTest        = errors.New("this receipt is from the test environment, but it was sent to the production environment for verification. Send it to the test environment instead")
	ErrReceiptIsForProduction  = errors.New("this receipt is from the production environment, but it was sent to the test environment for verification. Send it to the production environment instead")
	ErrReceiptUnauthorized     = errors.New("this receipt could not be authorized. Treat this the same as if a purchase was never made")
	ErrSubscriptionExpired     = errors.New("subscription has expired")
	ErrDuplicateReceipt        = errors.New("duplicate receipt")
	ErrInternalDataAccessError = errors.New("internal data access error")
	ErrUnknown                 = errors.New("an unknown error occurred")
)

// Returns error message by status code
func HandleErrors(status int) error {
	var e error
	switch status {
	case 0:
		return nil
	case 21000:
		e = ErrInvalidJSON
	case 21002:
		e = ErrInvalidReceiptData
	case 21003:
		e = ErrReceiptUnauthenticated
	case 21004:
		e = ErrInvalidSharedSecret
	case 21005:
		e = ErrServerUnavailable
	case 21006:
		e = ErrSubscriptionExpired
	case 21007:
		e = ErrReceiptIsForTest
	case 21008:
		e = ErrReceiptIsForProduction
	case 21009:
		e = ErrInternalDataAccessError
	case 21010:
		e = ErrReceiptUnauthorized
	default:
		if status >= 21100 && status <= 21199 {
			e = ErrInternalDataAccessError
		} else {
			e = ErrUnknown
		}
	}

	return e
}
