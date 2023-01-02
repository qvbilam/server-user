package jwt

import "github.com/dgrijalva/jwt-go"

type CustomClaims struct {
	ID       int64
	Code     int64
	Nickname string
	jwt.StandardClaims
}
