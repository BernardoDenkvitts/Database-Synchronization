package service

import (
	"github.com/BernardoDenkvitts/PostgresAPP/internal/infra"
	"github.com/BernardoDenkvitts/PostgresAPP/internal/types"
)

type UserService interface {
	CreateUser(types.UserRequestDTO) (string, error)
	GetUsers() ([]*types.UserResponseDTO, error)
	GetUserById(string) (*types.UserResponseDTO, error)
}

type UserServiceImpl struct {
	db infra.Storage
}

func NewUserServiceImpl(db infra.Storage) *UserServiceImpl {
	return &UserServiceImpl{
		db: db,
	}
}

func (svc *UserServiceImpl) CreateUser(userRequestDTO types.UserRequestDTO) (string, error) {
	user := types.NewUser(userRequestDTO)
	if err := svc.db.CreateUserInformation(user); err != nil {
		return "", err
	}
	return user.Id, nil
}

func (svc *UserServiceImpl) GetUsers() ([]*types.UserResponseDTO, error) {
	users, err := svc.db.GetUsersInformations()
	if err != nil {
		return nil, err
	}

	usersResponseDTO := make([]*types.UserResponseDTO, len(users))
	for i, user := range users {
		usersResponseDTO[i] = types.NewUserResponseDTO(*user)
	}

	return usersResponseDTO, nil
}

func (svc *UserServiceImpl) GetUserById(id string) (*types.UserResponseDTO, error) {
	user, err := svc.db.GetUserById(id)
	if err != nil {
		return nil, err
	}

	if user.Id == "" {
		return nil, nil
	}

	return types.NewUserResponseDTO(*user), nil
}
