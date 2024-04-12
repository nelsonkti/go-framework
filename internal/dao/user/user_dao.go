package user

import "fmt"

type UserDao struct {
}

func NewUserDao() *UserDao {
	return &UserDao{}
}

func (d *UserDao) QueryAll() {
	fmt.Println("UserDao QueryAll")
}
