package post

import (
	"context"
	"net/http"

	"github.com/forums/app/models"
)

type PostHandler interface {
	GetDetails(w http.ResponseWriter, r *http.Request)
	UpdateDetails(w http.ResponseWriter, r *http.Request)
	CreatePosts(w http.ResponseWriter, r *http.Request)
}

type PostRepo interface {
	GetPost(ctx context.Context, id int) (*models.Post, error)
	UpdateMessage(ctx context.Context, request *models.MessagePostRequest) error
	CreatePosts(ctx context.Context, posts *[]models.Post) (*[]models.Post, error)
	CreateForumsUsers(ctx context.Context, posts *[]models.Post) error
	GetPostsThread(ctx context.Context, id int) (int, error)
	ClearCache()
}
