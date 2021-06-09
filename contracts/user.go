package contracts

import "time"

type IUser interface {
	DeletedAt() *time.Time
}

type IUserService interface {
	FindByPhone(phone string) (IUser, error)
}
