# Go client library for signing in with apple

For more information about apple sign in, please review [apple doc](https://developer.apple.com/documentation/sign_in_with_apple/sign_in_with_apple_rest_api).

## Install

```bash
go get github.com/canopas/apple-sdk-go/auth
```

## How to use?

- **TeamId** :  10-char App Id prefix found in App identifiers section (Ex: AB65CD4321)

- **ClientId** : ClientID is the "Services ID" value that you get when navigating to your "sign in with Apple"-enabled service ID (Ex: com.example.me)

- **KeyId** : This is the ID of the private key (Ex: FE12DC34BA)

- **Secret** : This is the private key file (.p8). You can download it from [apple portal](https://developer.apple.com/account/resources/)


```go

// Create new secret request with default client
req, err := auth.WithDefaultClient("team-id", "client-id", "key-id", "secret-key-file-path")

if err != nil {
	log.Fatal(err.Error())
}

// OR
// Create new secret request with custom client
client := &http.Client{
	Timeout: 10 * time.Second,
}

req, err := auth.WithCustomClient(client, "team-id", "client-id", "key-id", "secret-key-file-path")

if err != nil {
	log.Fatal(err.Error())
}

// To do authorization request validation with authorization code from mobile app
resp, err := req.ValidateCode(context.Background(), "auth-code") 

// OR
// To do authorization request validation with authorization code and redirect uri from web app
resp, err := req.ValidateCodeWithRedirectURI(context.Background(), "auth-code", "redirect-uri") 

// OR
// Refresh token validation request
resp, err := req.ValidateRefreshToken(context.Background(), "refresh-token") 

if err != nil {
	log.Fatal(err.Error())
}

// get user
user, err := resp.GetUser()

if err != nil {
	log.Fatal(err.Error())
}

log.Println(user)

// get user's uniqueId
id, err := resp.UniqueID()

if err != nil {
	log.Fatal(err.Error())
}

log.Println(id)

// get user email
email, err := resp.Email()

if err != nil {
	log.Fatal(err.Error())
}

log.Println(email)

// Get user status 
// The possible values are: 0 (or Unsupported), 1 (or Unknown), 2 (or LikelyReal)
userStatus, err := resp.RealUserStatus()

if err != nil {
	log.Fatal(err.Error())
}

log.Println(userStatus)

```
