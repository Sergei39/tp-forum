package forum

import (
	"context"
	"net/http"

	"github.com/forums/app/models"
	"github.com/forums/utils/response"
)

type ForumHandler interface {
	CreateForum(w http.ResponseWriter, r *http.Request)
	GetDetails(w http.ResponseWriter, r *http.Request)
	GetUsers(w http.ResponseWriter, r *http.Request)
	GetThreads(w http.ResponseWriter, r *http.Request)
}

type ForumUsecase interface {
	CreateForum(ctx context.Context, forum *models.Forum) (*response.Response, error)
	GetUsers(ctx context.Context, forumUsers *models.ForumUsers) (*response.Response, error)
	GetThreads(ctx context.Context, forumThreads *models.ForumThreads) (*response.Response, error)
	GetForumBySlug(ctx context.Context, slug string) (*response.Response, error)
}

type ForumRepo interface {
	CreateForum(ctx context.Context, forum *models.Forum) (int, error)
	GetForumBySlug(ctx context.Context, title string) (*models.Forum, error)
	GetUsers(ctx context.Context, forumUsers *models.ForumUsers) (*[]models.User, error)
	GetThreads(ctx context.Context, forumThreads *models.ForumThreads) (*[]models.Thread, error)
}
