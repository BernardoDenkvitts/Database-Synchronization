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
}

type UserRouteImpl struct {
	userService service.UserService
}

func NewUserRoute(userService service.UserService) *UserRouteImpl {
	return &UserRouteImpl{
		userService: userService,
	}
}

func (userRoute *UserRouteImpl) Routes(router *http.ServeMux) {
	router.Handle("/user/", http.StripPrefix("/user", router))
	router.HandleFunc("POST /create", userRoute.handleCreateUser)
}

func (userRoute *UserRouteImpl) handleCreateUser(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("uri", "/mongodb/user/"+newUser.Id)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(types.ApiResponse{Status: http.StatusCreated, Response: "Created"})
}
