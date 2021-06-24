package delivery

import (
	"encoding/json"
	"net/http"

	userModel "github.com/forums/app/internal/user"
	"github.com/forums/app/models"
	"github.com/forums/utils/logger"
	"github.com/forums/utils/response"
	"github.com/gorilla/mux"
)

type handler struct {
	userRepo userModel.UserRepo
}

func NewUserHandler(userRepo userModel.UserRepo) userModel.UserHandler {
	return &handler{
		userRepo: userRepo,
	}
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	nickname := vars["nickname"]
	newUser := new(models.User)
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	newUser.Nickname = nickname
	logger.Delivery().Info(ctx, logger.Fields{"request data": *newUser})

	users, err := h.userRepo.GetUserByNameAndEmail(ctx, newUser.Nickname, newUser.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(*users) != 0 {
		response.New(http.StatusConflict, users).SendSuccess(w)
		return
	}

	err = h.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.New(http.StatusCreated, newUser).SendSuccess(w)
}

func (h *handler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	nickname := vars["nickname"]
	logger.Delivery().Info(ctx, logger.Fields{"request data": nickname})

	user, err := h.userRepo.GetUserByName(ctx, nickname)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if user == nil {
		message := models.Message{
			Message: "Can't find user with id #" + nickname + "\n",
		}
		response.New(http.StatusNotFound, message).SendSuccess(w)
		return
	}

	response.New(http.StatusOK, user).SendSuccess(w)
}

func (u *handler) fixData(newUser *models.User, oldUser *models.User) *models.User {
	if newUser.About == "" {
		newUser.About = oldUser.About
	}

	if newUser.Email == "" {
		newUser.Email = oldUser.Email
	}

	if newUser.Fullname == "" {
		newUser.Fullname = oldUser.Fullname
	}

	return newUser
}

func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	nickname := vars["nickname"]
	newUser := new(models.User)
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	newUser.Nickname = nickname
	logger.Delivery().Info(ctx, logger.Fields{"request data": *newUser})

	userDb, err := h.userRepo.GetUserByName(ctx, newUser.Nickname)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if userDb == nil {
		message := models.Message{
			Message: "Can't find user with id #" + newUser.Nickname + "\n",
		}
		response.New(http.StatusNotFound, message).SendSuccess(w)
		return
	}

	newUser = h.fixData(newUser, userDb)

	userDb, err = h.userRepo.GetUserByEmail(ctx, newUser.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if userDb != nil && userDb.Nickname != newUser.Nickname {
		response.New(http.StatusConflict, userDb).SendSuccess(w)
		return
	}

	_, err = h.userRepo.UpdateUser(ctx, newUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.New(http.StatusOK, newUser).SendSuccess(w)
}
