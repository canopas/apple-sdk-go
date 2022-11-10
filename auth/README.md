# Go client library for signing in user

Go client library for signing in user.

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

req := auth.New("team-id", "client-id", "key-id", "secret-key")

resp, err := req.ValidateCode("auth-code") 

// or

resp, err := req.ValidateCodeWithRedirectURI("auth-code", "redirect-uri") 

// or

resp, err := req.ValidateRefreshToken("refresh-token") 

if err != nil {
	log.Fatal(err.Error())
}

// get validated user
user, err := resp.GetUser(resp.IDToken)

if err != nil {
	log.Fatal(err.Error())
}

log.Println(user)

// get validated user's uniqueId
id, err := resp.UniqueID(resp.IDToken)

if err != nil {
	log.Fatal(err.Error())
}

log.Println(id)

// get validated user email
email, err := resp.Email(resp.IDToken)

if err != nil {
	log.Fatal(err.Error())
}

log.Println(email)

// get validated user status
userStatus, err := resp.RealUserStatus(resp.IDToken)

if err != nil {
	log.Fatal(err.Error())
}

log.Println(userStatus)

```
