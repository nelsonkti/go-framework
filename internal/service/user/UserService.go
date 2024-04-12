package user

import (
	"go-framework/internal/dao/user"
)

type UserService struct {
	dao *user.UserDao
}

func NewUserService(dao *user.UserDao) *UserService {
	return &UserService{dao: dao}
}

func (s *UserService) Get() {
	s.dao.QueryAll()
}
