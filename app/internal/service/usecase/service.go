package usecase

import (
	"context"
	"net/http"

	postModel "github.com/forums/app/internal/post"
	serviceModel "github.com/forums/app/internal/service"
	"github.com/forums/utils/response"
)

type usecase struct {
	serviceRepo serviceModel.ServiceRepo
	postRepo    postModel.PostRepo
}

func NewServiceUsecase(serviceRepo serviceModel.ServiceRepo, postRepo postModel.PostRepo) serviceModel.ServiceUsecase {
	return &usecase{
		serviceRepo: serviceRepo,
		postRepo:    postRepo,
	}
}

func (u *usecase) ClearDb(ctx context.Context) error {
	u.postRepo.ClearCache()
	return u.serviceRepo.ClearDb(ctx)
}

func (u *usecase) StatusDb(ctx context.Context) (response.Response, error) {
	result, err := u.serviceRepo.StatusDb(ctx)
	if err != nil {
		return nil, err
	}

	response := response.New(http.StatusOK, result)
	return response, nil
}
