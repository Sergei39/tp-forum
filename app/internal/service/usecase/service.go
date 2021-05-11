package usecase

import (
	"context"
	"net/http"

	serviceModel "github.com/forums/app/internal/service"
	"github.com/forums/utils/response"
)

type usecase struct {
	serviceRepo serviceModel.ServiceRepo
}

func NewServiceUsecase(serviceRepo serviceModel.ServiceRepo) serviceModel.ServiceUsecase {
	return &usecase{
		serviceRepo: serviceRepo,
	}
}

func (u *usecase) ClearDb(ctx context.Context) error {
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
