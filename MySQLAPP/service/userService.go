package service

import (
	"github.com/BernardoDenkvitts/MySQLApp/storage"
	"github.com/BernardoDenkvitts/MySQLApp/types"
)

type UserService interface {
	CreateUser(userRequestDTO types.UserRequestDTO) error
	GetUsers() ([]*types.User, error)
}

type UserServiceImpl struct {
	storage storage.Storage
}

func NewUserService(storage storage.Storage) *UserServiceImpl {
	return &UserServiceImpl{
		storage: storage,
	}
}
