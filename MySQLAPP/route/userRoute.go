package route

import (
	"encoding/json"
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
	router.HandleFunc("GET /user", userRoute.handleGetUsersInformation)
	router.HandleFunc("GET /userteste", userRoute.Teste)
}

func (userRoute *UserRoute) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	userRequestDTO := new(types.UserRequestDTO)
	if err := utils.ParseJson(r, userRequestDTO); err != nil {
		utils.WriteJson(w, http.StatusBadRequest, types.ApiResponse{Status: http.StatusBadRequest, Response: err.Error()})
		return
	}

	newUser, err := userRoute.userService.CreateUser(*userRequestDTO)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, types.ApiResponse{Status: http.StatusInternalServerError, Response: err.Error()})
		return
	}

	w.Header().Add("Content-Type", "application-json")
	w.Header().Add("uri", "/user/get/"+newUser.Id)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(types.ApiResponse{Status: http.StatusCreated, Response: "Created"})
}

func (userRoute *UserRoute) handleGetUserInformationsById(w http.ResponseWriter, r *http.Request) {
	userResponseDTO, err := userRoute.userService.GetUserById(r.PathValue("id"))
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, types.ApiResponse{Status: http.StatusInternalServerError, Response: err.Error()})
		return
	}
	if userResponseDTO == nil {
		utils.WriteJson(w, http.StatusNotFound, types.ApiResponse{Status: http.StatusNotFound, Response: "User Not Found"})
		return
	}

	utils.WriteJson(w, http.StatusOK, types.ApiResponse{Status: http.StatusOK, Response: userResponseDTO})
}

func (userRoute *UserRoute) handleGetUsersInformation(w http.ResponseWriter, r *http.Request) {
	usersResponseDTO, err := userRoute.userService.GetUsers()
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, types.ApiResponse{Status: http.StatusInternalServerError, Response: err.Error()})
		return
	}

	utils.WriteJson(w, http.StatusOK, types.ApiResponse{Status: http.StatusOK, Response: usersResponseDTO})
}

func (userRoute *UserRoute) Teste(w http.ResponseWriter, r *http.Request) {
	usersResponseDTO, err := userRoute.userService.Teste()
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, types.ApiResponse{Status: http.StatusInternalServerError, Response: err.Error()})
		return
	}

	utils.WriteJson(w, http.StatusOK, types.ApiResponse{Status: http.StatusOK, Response: usersResponseDTO})
}
