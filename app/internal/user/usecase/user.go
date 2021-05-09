package usecase

import (
	"context"
	"net/http"

	userModel "github.com/forums/app/internal/user"
	"github.com/forums/app/models"
	"github.com/forums/utils/response"
)

type usecase struct {
	userRepo userModel.UserRepo
}

func NewUserUsecase(userRepo userModel.UserRepo) userModel.UserUsecase {
	return &usecase{
		userRepo: userRepo,
	}
}

func (u *usecase) CreateUser(ctx context.Context, user models.User) (
	response.Response, error) {

	userDb, err := u.userRepo.GetUserByName(ctx, user.Nickname)
	if err != nil {
		return nil, err
	}
	if userDb != nil {
		response := response.New(http.StatusConflict, user)
		return response, nil
	}

	_, err = u.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	response := response.New(http.StatusCreated, user)
	return response, nil
}

func (u *usecase) GetUserByName(ctx context.Context, name string) (
	response.Response, error) {

	user, err := u.userRepo.GetUserByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if user == nil {
		message := models.Message{
			Message: "Can't find user with id #" + name + "\n",
		}
		response := response.New(http.StatusNotFound, message)
		return response, nil
	}

	response := response.New(http.StatusOK, user)
	return response, nil
}
