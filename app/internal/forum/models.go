package forum

import (
	"context"
	"net/http"

	"github.com/forums/app/models"
)

type ForumHandler interface {
	CreateForum(w http.ResponseWriter, r *http.Request)
	GetDetails(w http.ResponseWriter, r *http.Request)
	GetUsers(w http.ResponseWriter, r *http.Request)
	GetThreads(w http.ResponseWriter, r *http.Request)
}

type ForumRepo interface {
	CreateForum(ctx context.Context, forum *models.Forum) (int, error)
	GetForumBySlug(ctx context.Context, title string) (*models.Forum, error)
	GetUsers(ctx context.Context, forumUsers *models.ForumUsers) (*[]models.User, error)
	GetThreads(ctx context.Context, forumThreads *models.ForumThreads) (*[]models.Thread, error)
}
