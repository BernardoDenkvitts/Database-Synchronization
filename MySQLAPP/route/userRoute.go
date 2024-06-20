package route

import (
	"fmt"
	"net/http"

	"github.com/BernardoDenkvitts/MySQLApp/service"
	"github.com/BernardoDenkvitts/MySQLApp/types"
	"github.com/BernardoDenkvitts/MySQLApp/utils"
)

type UserRoute struct {
	userService service.UserServiceImpl
}

func NewUserRoute(userService service.UserServiceImpl) *UserRoute {
	return &UserRoute{userService: userService}
}

func (userRoute *UserRoute) Routes(router *http.ServeMux) {
	router.Handle("/user/", http.StripPrefix("/user", router))
	router.HandleFunc("POST /create", userRoute.handleCreateUser)
	router.HandleFunc("GET /teste", userRoute.teste)
}

func (UserRoute *UserRoute) teste(w http.ResponseWriter, r *http.Request) {
	utils.WriteJson(w, http.StatusCreated, "JESUS CRISTO", nil)
}

func (UserRoute *UserRoute) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var payload types.UserRequestDTO
	if err := utils.ParseJson(r, payload); err != nil {
		err := utils.WriteJson(w, http.StatusBadRequest, err, nil)
		if err != nil {
			return
		}
	}
	newUser, err := UserRoute.userService.CreateUser(payload)
	if err != nil {
		err := utils.WriteJson(w, http.StatusInternalServerError, err, nil)
		if err != nil {
			return
		}
		return
	}
	fmt.Println(*newUser)
	fmt.Println(newUser.Id)
	var header = map[string]string{"location": "/user/get/" + newUser.Id}
	utils.WriteJson(w, http.StatusCreated, types.NewUserResponseDTO{Status: http.StatusCreated}, &header)
}
