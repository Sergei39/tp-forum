package user

import (
	"context"
	"net/http"

	"github.com/forums/app/models"
)

type UserHandler interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
}

type UserRepo interface {
	CreateUser(ctx context.Context, user *models.User) (err error)
	GetUserByName(ctx context.Context, name string) (*models.User, error)   // TODO: возможно удлаить метод
	GetUserByEmail(ctx context.Context, email string) (*models.User, error) // TODO: возможно удлаить метод
	UpdateUser(ctx context.Context, user *models.User) (id int, err error)
	GetUserByNameAndEmail(ctx context.Context, name, email string) (*[]models.User, error)
}
