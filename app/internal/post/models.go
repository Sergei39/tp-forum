package post

import (
	"context"
	"net/http"

	"github.com/forums/app/models"
	"github.com/forums/utils/response"
)

type PostHandler interface {
	GetDetails(w http.ResponseWriter, r *http.Request)
	UpdateDetails(w http.ResponseWriter, r *http.Request)
	CreatePosts(w http.ResponseWriter, r *http.Request)
}

type PostUsecase interface {
	GetDetails(ctx context.Context, request *models.RequestPost) (*response.Response, error)
	UpdateMessage(ctx context.Context, request *models.MessagePostRequest) (*response.Response, error)
	CreatePosts(ctx context.Context, posts *[]models.Post, slugOrId string) (*response.Response, error)
}

type PostRepo interface {
	GetPost(ctx context.Context, id int) (*models.Post, error)
	UpdateMessage(ctx context.Context, request *models.MessagePostRequest) error
	CreatePosts(ctx context.Context, posts *[]models.Post) (*[]models.Post, error)
	CreateForumsUsers(ctx context.Context, posts *[]models.Post) error
	GetPostsThread(ctx context.Context, id int) (int, error)
	ClearCache()
}
