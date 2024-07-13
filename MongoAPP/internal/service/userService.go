package service

import (
	"github.com/BernardoDenkvitts/MongoAPP/internal/infra"
	"github.com/BernardoDenkvitts/MongoAPP/internal/types"
)

type UserService interface {
	CreateUser(userRequestDTO types.UserRequestDTO) (string, error)
	GetUsers() ([]*types.UserResponseDTO, error)
	GetUserById(id string) (*types.UserResponseDTO, error)
}

type UserServiceImpl struct {
	Db infra.Storage
}

func NewUserServiceImpl(database infra.Storage) *UserServiceImpl {
	return &UserServiceImpl{
		Db: database,
	}
}

func (userService *UserServiceImpl) CreateUser(userRequestDTO types.UserRequestDTO) (string, error) {
	newUser := types.NewUser(userRequestDTO)
	err := userService.Db.CreateUserInformation(newUser)
	if err != nil {
		return "", err
	}

	return newUser.Id, nil
}

func (userService *UserServiceImpl) GetUserById(id string) (*types.UserResponseDTO, error) {
	user, err := userService.Db.GetUserById(id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	return types.NewUserResponseDTO(*user), nil
}

func (userService *UserServiceImpl) GetUsers() ([]*types.UserResponseDTO, error) {
	users, err := userService.Db.GetUsersInformations()
	if err != nil {
		return nil, err
	}

	usersResponseDTO := make([]*types.UserResponseDTO, len(users))
	for i, user := range users {
		usersResponseDTO[i] = types.NewUserResponseDTO(*user)
	}

	return usersResponseDTO, nil
}
