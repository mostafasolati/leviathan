package models

import (
	"github.com/dgrijalva/jwt-go"
)

// UserClaims contains jwt data
type UserClaims struct {
	jwt.StandardClaims
	ID    int      `json:"id"`
	Name  string   `json:"name"`
	Roles []string `json:"roles"`
	Phone string   `json:"phone"`
}

// Token is an object contain access token and refresh token for authenticating via jwt
type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
