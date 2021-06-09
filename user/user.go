package user

import "github.com/mostafasolati/leviathan/contracts"

type user struct{}

func NewUserService() contracts.IUserService {
	return &user{}
}

func (s *user) FindByPhone(phone string) (contracts.IUser, error) {
	panic("implement me")
}
