package model

import "github.com/golang-jwt/jwt/v4"

type JwtClaims struct {
	Uid string `json:"uid"`
	jwt.RegisteredClaims
}
