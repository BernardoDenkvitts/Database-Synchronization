package service

import (
	"github.com/BernardoDenkvitts/MySQLApp/internal/infra"
	"github.com/BernardoDenkvitts/MySQLApp/internal/types"
)

type UserService interface {
	CreateUser(userRequestDTO types.UserRequestDTO) error
	GetUsers() ([]*types.UserResponseDTO, error)
	GetUserById(id string) (*types.UserResponseDTO, error)
}

type UserServiceImpl struct {
	storage infra.Storage
}

func NewUserService(storage infra.Storage) *UserServiceImpl {
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

func (userService *UserServiceImpl) GetUsers() ([]*types.UserResponseDTO, error) {
	users, err := userService.storage.GetUsersInformations()
	if err != nil {
		return nil, err
	}

	usersResponseDTO := make([]*types.UserResponseDTO, len(users))
	for idx, user := range users {
		usersResponseDTO[idx] = types.NewUserResponseDTO(*user)
	}

	return usersResponseDTO, nil
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
