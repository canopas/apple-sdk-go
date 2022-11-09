# Generate Secret for authentication API

Go client library for generating secret

## Install

```bash
go get github.com/canopas/apple-sdk-go/secret
```

## How to use?

- **TeamId** :  10-char App Id prefix found in App identifiers section (Ex: AB65CD4321)

- **ClientId** : ClientID is the "Services ID" value that you get when navigating to your "sign in with Apple"-enabled service ID (Ex: com.example.me)

- **KeyId** : This is the ID of the private key (Ex: FE12DC34BA)

- **Secret** : This is the private key file (.p8). You can download it from [apple portal](https://developer.apple.com/account/resources/)

```go
    req := secret.New("teamId", "clientId", "keyId", "secret")

    secret, err := req.GenerateClientSecret()

    if err != nil {
        log.Fatal(err)
    }
```

