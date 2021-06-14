package delivery

import (
	"encoding/json"
	"net/http"

	userModel "github.com/forums/app/internal/user"
	"github.com/forums/app/models"
	"github.com/forums/utils/errors"
	"github.com/forums/utils/logger"
	"github.com/gorilla/mux"
)

type handler struct {
	userUsecase userModel.UserUsecase
}

func NewUserHandler(usecase userModel.UserUsecase) userModel.UserHandler {
	return &handler{
		userUsecase: usecase,
	}
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	nickname := vars["nickname"]
	newUser := new(models.User)
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		sendErr := errors.New(http.StatusBadRequest, err.Error())
		logger.Delivery().Error(ctx, sendErr)
		w.WriteHeader(sendErr.Code())
		return
	}
	defer r.Body.Close()

	newUser.Nickname = nickname
	logger.Delivery().Info(ctx, logger.Fields{"request data": *newUser})

	response, err := h.userUsecase.CreateUser(ctx, *newUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w)
}

func (h *handler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	nickname := vars["nickname"]
	logger.Delivery().Info(ctx, logger.Fields{"request data": nickname})

	response, err := h.userUsecase.GetUserByName(ctx, nickname)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w)
}

func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	nickname := vars["nickname"]
	newUser := new(models.User)
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		sendErr := errors.New(http.StatusBadRequest, err.Error())
		logger.Delivery().Error(ctx, sendErr)
		w.WriteHeader(sendErr.Code())
		return
	}
	defer r.Body.Close()

	newUser.Nickname = nickname
	logger.Delivery().Info(ctx, logger.Fields{"request data": *newUser})

	response, err := h.userUsecase.UpdateUser(ctx, *newUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w)
}
