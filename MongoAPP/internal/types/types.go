package types

import (
	"time"

	"github.com/google/uuid"
)

type ApiResponse struct {
	Status   int `json:"status"`
	Response any `json:"response"`
}

type User struct {
	Id        string    `bson:"id"`
	FirstName string    `bson:"firstName"`
	LastName  string    `bson:"lastName"`
	CreatedAt time.Time `bson:"created_at"`
}

func NewUser(userRequestDTO UserRequestDTO) *User {
	return &User{
		Id:        uuid.Must(uuid.NewRandom()).String(),
		FirstName: userRequestDTO.FirstName,
		LastName:  userRequestDTO.LastName,
		CreatedAt: time.Now().UTC(),
	}
}

type UserRequestDTO struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type UserResponseDTO struct {
	Id        string    `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewUserResponseDTO(user User) *UserResponseDTO {
	return &UserResponseDTO{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
	}
}
