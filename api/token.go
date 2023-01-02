package api

import (
	"github.com/dgrijalva/jwt-go"
	"time"
	proto "user/api/qvbilam/user/v1"
	AuthJWT "user/auth/jwt"
	"user/global"
	"user/middleware"
)

func GenerateUserToken(entity *proto.UserResponse) string {
	j := middleware.NewJWT()
	expire := time.Now().Unix() + global.ServerConfig.JWTConfig.Expire
	issuer := global.ServerConfig.JWTConfig.Issuer
	claims := AuthJWT.CustomClaims{
		ID:       entity.Id,
		Code:     entity.Code,
		Nickname: entity.Nickname,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expire,
			Issuer:    issuer,
			NotBefore: time.Now().Unix(),
		},
	}

	token, err := j.CreateToken(claims)
	if err != nil {
		return ""
	}

	return token
}
