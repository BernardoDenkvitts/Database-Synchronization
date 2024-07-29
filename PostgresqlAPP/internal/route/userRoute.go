package route

import (
	"encoding/json"
	"net/http"

	"github.com/BernardoDenkvitts/PostgresqlAPP/internal/service"
	"github.com/BernardoDenkvitts/PostgresqlAPP/internal/types"
	"github.com/BernardoDenkvitts/PostgresqlAPP/internal/utils"
)

type UserRoutes interface {
	Routes(router *http.ServeMux)
	handleCreateUser(w http.ResponseWriter, r *http.Request)
	handleGetUserInformationsById(w http.ResponseWriter, r *http.Request)
	handleGetUsersInformation(w http.ResponseWriter, r *http.Request)
}

type UserRoutesImpl struct {
	userService service.UserService
}

func NewUserRoutesImpl(userService service.UserService) *UserRoutesImpl {
	return &UserRoutesImpl{
		userService: userService,
	}
}

func (route *UserRoutesImpl) Routes(router *http.ServeMux) {
	router.Handle("/user/", http.StripPrefix("/user", router))
	router.HandleFunc("POST /create", route.handleCreateUser)
	router.HandleFunc("GET /user/{id}", route.handleGetUserInformationsById)
	router.HandleFunc("GET /user", route.handleGetUsersInformation)
}

func (route *UserRoutesImpl) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	userRequestDTO := new(types.UserRequestDTO)
	if err := utils.ParseJson(r, userRequestDTO); err != nil {
		utils.WriteJson(w, http.StatusBadRequest, types.ApiResponse{Status: http.StatusBadRequest, Response: err.Error()})
		return
	}

	userId, err := route.userService.CreateUser(*userRequestDTO)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, types.ApiResponse{Status: http.StatusInternalServerError, Response: err.Error()})
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("uri", "/postgres/user/"+userId)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(types.ApiResponse{Status: http.StatusCreated, Response: "Created"})
}

func (route *UserRoutesImpl) handleGetUserInformationsById(w http.ResponseWriter, r *http.Request) {
	userResponseDTO, err := route.userService.GetUserById(r.PathValue("id"))
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

func (route *UserRoutesImpl) handleGetUsersInformation(w http.ResponseWriter, r *http.Request) {
	usersResponseDTO, err := route.userService.GetUsers()
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, types.ApiResponse{Status: http.StatusInternalServerError, Response: err.Error()})
		return
	}

	utils.WriteJson(w, http.StatusOK, types.ApiResponse{Status: http.StatusOK, Response: usersResponseDTO})
}
