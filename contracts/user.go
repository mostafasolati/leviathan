package contracts

import "time"

type IUser interface {
	ID() int
	DeletedAt() *time.Time
	FullName() string
	FirstName() string
	LastName() string
	Phone() string
	RefreshTokenExpiry() *time.Time
	RefreshToken() string
	SetRefreshToken(token string)
	SetRefreshTokenExpiry(t time.Time)
}

type IUserService interface {
	FindByPhone(phone string) (IUser, error)
	FindByRefreshToken(token string) (IUser, error)
	Create() error
}
