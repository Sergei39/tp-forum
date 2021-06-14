package service

import (
	"context"
	"net/http"

	"github.com/forums/app/models"
	"github.com/forums/utils/response"
)

type ServiceHandler interface {
	ClearDb(w http.ResponseWriter, r *http.Request)
	StatusDb(w http.ResponseWriter, r *http.Request)
}

type ServiceUsecase interface {
	ClearDb(ctx context.Context) error
	StatusDb(ctx context.Context) (response.Response, error)
}

type ServiceRepo interface {
	ClearDb(ctx context.Context) error
	StatusDb(ctx context.Context) (*models.InfoStatus, error)
}
