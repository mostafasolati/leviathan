package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// TypeRole is a custom data type for representing users role
type TypeRole uint8

// We can add new roles in between them
const (
	TypeRoleGuest         TypeRole = 10
	TypeRoleUser          TypeRole = 20
	TypeRoleCourier       TypeRole = 30
	TypeRoleAmbassador    TypeRole = 40
	TypeRoleTiendaManager TypeRole = 41
	TypeRoleSupport       TypeRole = 50
	TypeRoleAdmin         TypeRole = 60
)

var roleNames = map[TypeRole]string{
	TypeRoleGuest:         "guest",
	TypeRoleUser:          "user",
	TypeRoleCourier:       "courier",
	TypeRoleAmbassador:    "ambassador",
	TypeRoleSupport:       "support",
	TypeRoleAdmin:         "admin",
	TypeRoleTiendaManager: "tienda_manager",
}

// Name returns the role name.
func (r TypeRole) Name() string {
	if name, ok := roleNames[r]; ok {
		return name
	}
	return "unknown"
}

// TypeRoleForName returns the role matching a name.
func TypeRoleForName(name string) TypeRole {
	for role, roleName := range roleNames {
		if roleName == name {
			return role
		}
	}
	return TypeRoleUser
}

// User is the data model representing a user in the system
type User struct {
	Base
	FirstName          string    `json:"first_name"`
	LastName           string    `json:"last_name"`
	Phone              string    `json:"phone"`
	Email              string    `json:"email"`
	RefreshToken       string    `json:"-"`
	RefreshTokenExpiry time.Time `json:"-"`
}

// FullName concatenates first name and last name
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// UserClaims contains jwt data
type UserClaims struct {
	jwt.StandardClaims
	ID    int      `json:"id"`
	Name  string   `json:"name"`
	Roles []string `json:"roles"`
	Phone string   `json:"phone"`
}

// HasRole returns true if the claims have a specific role.
func (c *UserClaims) HasRole(role TypeRole) bool {
	roleName := role.Name()
	for _, r := range c.Roles {
		if r == roleName {
			return true
		}
	}
	return false
}

// Token is an object contain access token and refresh token for authenticating via jwt
type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
