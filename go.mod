module github.com/canopas/apple-sdk-go

go 1.18

replace auth => ./auth

require auth v0.0.0-00010101000000-000000000000

require (
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/golang-jwt/jwt/v4 v4.4.2 // indirect
)
