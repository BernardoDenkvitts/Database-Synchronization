package route

import (
	"encoding/json"
	"net/http"

	"github.com/BernardoDenkvitts/MongoAPP/internal/service"
	"github.com/BernardoDenkvitts/MongoAPP/internal/types"
	"github.com/BernardoDenkvitts/MongoAPP/internal/utils"
)

type UserRoutes interface {
	Routes(router *http.ServeMux)
	handleCreateUser(w http.ResponseWriter, r *http.Request)
	handleGetUserInformationsById(w http.ResponseWriter, r *http.Request)
	handleGetUsersInformation(w http.ResponseWriter, r *http.Request)
	teste(w http.ResponseWriter, r *http.Request)
}

type UserRouteImpl struct {
	userService service.UserService
}

func NewUserRouteImpl(userService service.UserService) *UserRouteImpl {
	return &UserRouteImpl{
		userService: userService,
	}
}

func (userRoute *UserRouteImpl) Routes(router *http.ServeMux) {
	router.Handle("/user/", http.StripPrefix("/user", router))
	router.HandleFunc("POST /create", userRoute.handleCreateUser)
	router.HandleFunc("GET /user/{id}", userRoute.handleGetUserInformationsById)
	router.HandleFunc("GET /user", userRoute.handleGetUsersInformation)
	router.HandleFunc("GET /user/teste", userRoute.teste)
}

func (userRoute *UserRouteImpl) teste(w http.ResponseWriter, r *http.Request) {
	usersResponseDTO, err := userRoute.userService.Teste()
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, types.ApiResponse{Status: http.StatusInternalServerError, Response: err.Error()})
		return
	}

	utils.WriteJson(w, http.StatusOK, types.ApiResponse{Status: http.StatusOK, Response: usersResponseDTO})
}

func (userRoute *UserRouteImpl) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	userRequestDTO := new(types.UserRequestDTO)
	if err := utils.ParseJson(r, userRequestDTO); err != nil {
		utils.WriteJson(w, http.StatusBadRequest, types.ApiResponse{Status: http.StatusBadRequest, Response: err.Error()})
		return
	}

	newUserID, err := userRoute.userService.CreateUser(*userRequestDTO)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, types.ApiResponse{Status: http.StatusInternalServerError, Response: err.Error()})
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("uri", "/mongodb/user/"+newUserID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(types.ApiResponse{Status: http.StatusCreated, Response: "Created"})
}

func (userRoute *UserRouteImpl) handleGetUserInformationsById(w http.ResponseWriter, r *http.Request) {
	user, err := userRoute.userService.GetUserById(r.PathValue("id"))
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, types.ApiResponse{Status: http.StatusInternalServerError, Response: err.Error()})
		return
	}

	if user == nil {
		utils.WriteJson(w, http.StatusNotFound, types.ApiResponse{Status: http.StatusNotFound, Response: "User Not Found"})
		return
	}

	utils.WriteJson(w, http.StatusOK, types.ApiResponse{Status: http.StatusOK, Response: user})
}

func (userRoute UserRouteImpl) handleGetUsersInformation(w http.ResponseWriter, r *http.Request) {
	usersResponseDTO, err := userRoute.userService.GetUsers()
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, types.ApiResponse{Status: http.StatusInternalServerError, Response: err.Error()})
		return
	}

	utils.WriteJson(w, http.StatusOK, types.ApiResponse{Status: http.StatusOK, Response: usersResponseDTO})
}
