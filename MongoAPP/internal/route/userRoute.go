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
	router.HandleFunc("GET /user/teste", route.teste)
}

func (route *UserRoutesImpl) teste(w http.ResponseWriter, r *http.Request) {
	usersResponseDTO, err := route.userService.Teste()
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, types.ApiResponse{Status: http.StatusInternalServerError, Response: err.Error()})
		return
	}

	utils.WriteJson(w, http.StatusOK, types.ApiResponse{Status: http.StatusOK, Response: usersResponseDTO})
}

func (route *UserRoutesImpl) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	userRequestDTO := new(types.UserRequestDTO)
	if err := utils.ParseJson(r, userRequestDTO); err != nil {
		utils.WriteJson(w, http.StatusBadRequest, types.ApiResponse{Status: http.StatusBadRequest, Response: err.Error()})
		return
	}

	newUserID, err := route.userService.CreateUser(*userRequestDTO)
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, types.ApiResponse{Status: http.StatusInternalServerError, Response: err.Error()})
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("uri", "/mongodb/user/"+newUserID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(types.ApiResponse{Status: http.StatusCreated, Response: "Created"})
}

func (route *UserRoutesImpl) handleGetUserInformationsById(w http.ResponseWriter, r *http.Request) {
	user, err := route.userService.GetUserById(r.PathValue("id"))
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

func (route UserRoutesImpl) handleGetUsersInformation(w http.ResponseWriter, r *http.Request) {
	usersResponseDTO, err := route.userService.GetUsers()
	if err != nil {
		utils.WriteJson(w, http.StatusInternalServerError, types.ApiResponse{Status: http.StatusInternalServerError, Response: err.Error()})
		return
	}

	utils.WriteJson(w, http.StatusOK, types.ApiResponse{Status: http.StatusOK, Response: usersResponseDTO})
}
