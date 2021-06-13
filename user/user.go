package user

import "github.com/mostafasolati/leviathan/contracts"

type user struct{}

func NewUserService() contracts.IUserService {
	return &user{}
}

func (s *user) FindByRefreshToken(token string) (contracts.IUser, error) {
	panic("implement me")
}

func (s *user) Create() error {
	panic("implement me")
}

func (s *user) FindByPhone(phone string) (contracts.IUser, error) {
	panic("implement me")
}
