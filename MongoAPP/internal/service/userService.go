package service

import "github.com/BernardoDenkvitts/MongoAPP/internal/types"

type UserService interface {
	CreateUser(userRequestDTO types.UserRequestDTO) (*types.User, error)
}

type UserServiceImpl struct{}

func NewUserService() *UserServiceImpl {
	return &UserServiceImpl{}
}

func (userService *UserServiceImpl) CreateUser(userRequestDTO types.UserRequestDTO) (*types.User, error) {
	newUser := types.NewUser(userRequestDTO)
	return newUser, nil
}
