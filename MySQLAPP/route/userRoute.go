package route

import (
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
	router.HandleFunc("GET /user/{id}", userRoute.handleGetUserInformationsById)
	router.HandleFunc("GET /user", userR)
}

func (userRoute *UserRoute) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	userRequestDTO := new(types.UserRequestDTO)
	if err := utils.ParseJson(r, userRequestDTO); err != nil {
		utils.WriteJson(w, http.StatusBadRequest, types.ApiResponse{Status: http.StatusBadRequest, Message: err.Error()}, nil)
		return
	}
	newUser, err := userRoute.userService.CreateUser(*userRequestDTO)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, types.ApiResponse{Status: http.StatusInternalServerError, Message: err.Error()}, nil)
		return
	}
	header := map[string]string{"url": "/user/get/" + newUser.Id}
	utils.WriteJson(w, http.StatusCreated, types.ApiResponse{Status: http.StatusCreated, Message: "Created"}, &header)
}

func (userRoute *UserRoute) handleGetUserInformationsById(w http.ResponseWriter, r *http.Request) {
	userResponseDTO, err := userRoute.userService.GetUserById(r.PathValue("id"))
	if err != nil {
		utils.WriteJson(w, http.StatusNotFound, types.ApiResponse{Status: http.StatusInternalServerError, Message: err.Error()}, nil)
		return
	}
	if userResponseDTO == nil {
		utils.WriteJson(w, http.StatusNotFound, types.ApiResponse{Status: http.StatusNotFound, Message: "User Not Found"}, nil)
		return
	}

	utils.WriteJson(w, http.StatusNotFound, types.ApiResponse{Status: http.StatusOK, Message: userResponseDTO}, nil)
}
