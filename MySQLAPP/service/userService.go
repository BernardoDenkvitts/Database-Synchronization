package service

import (
	"github.com/BernardoDenkvitts/MySQLApp/storage"
	"github.com/BernardoDenkvitts/MySQLApp/types"
)

type UserService interface {
	CreateUser(userRequestDTO types.UserRequestDTO) error
	GetUsers() ([]*types.UserResponseDTO, error)
	GetUserById(id string) (*types.UserResponseDTO, error)
}

type UserServiceImpl struct {
	storage storage.Storage
}

func NewUserService(storage storage.Storage) *UserServiceImpl {
	return &UserServiceImpl{
		storage: storage,
	}
}

func (userService *UserServiceImpl) CreateUser(userRequestDTO types.UserRequestDTO) (*types.User, error) {
	newUser := types.NewUser(userRequestDTO)
	err := userService.storage.CreateUserInformation(newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (userService *UserServiceImpl) GetUserById(id string) (*types.UserResponseDTO, error) {
	user, err := userService.storage.GetUserById(id)
	if err != nil {
		return nil, err
	}

	if user.Id == "" {
		return nil, nil
	}
	return &types.UserResponseDTO{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
	}, nil
}
