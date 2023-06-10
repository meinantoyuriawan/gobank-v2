package config

import "github.com/golang-jwt/jwt/v4"

var JWT_KEY = []byte("udbsacoi991283742gga987912g")

type JWTClaim struct {
	Username string
	jwt.RegisteredClaims
}
