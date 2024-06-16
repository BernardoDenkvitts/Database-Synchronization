package types

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        string
	FirstName string
	LastName  string
	CreatedAt time.Time
}

type UserRequestDTO struct {
	FirstName string `json: "firstName"`
	LastName  string `json: "lastName"`
}

func NewUser(userRequestDTO UserRequestDTO) *User {
	return &User{
		Id:        uuid.Must(uuid.NewRandom()).String(),
		FirstName: userRequestDTO.FirstName,
		LastName:  userRequestDTO.LastName,
		CreatedAt: time.Now().UTC(),
	}
}
