# Go client library for verifying receipt-data from appstore

Go client library for verifying receipt-data from appstore.

For more information about appstore receipts, please review [apple doc](https://developer.apple.com/documentation/appstorereceipts).

## Install

```bash
go get github.com/canopas/apple-sdk-go/receipt
```

## How to use?

- **ReceiptData** :  base64 encoded receipt data

- **Password** : Appstore shared secret (Ex: 8h1b362b673b442d8ali4b2141912948)

```go

// Create new IAP request with default client
client := receipt.New()

// OR
// Create new IAP request with custom client
httpCli := &http.Client{
	Timeout: 10 * time.Second,
}

client := receipt.NewWithClient(httpCli)

// verify receipt data
response, err := client.Verify(context.Background(), receipt.IAPRequest{
	ReceiptData: "your-receipt-data",
	Password:    "shared-secret",
})

if err != nil {
	log.Fatal(err.Error())
}

// handle errors if any
if response.Status != 0 {
	err = receipt.HandleErrors(response.Status)
	log.Fatal(err.Error())
}

log.Println(response)

```
