package thread

import (
	"context"
	"net/http"

	"github.com/forums/app/models"
)

type ThreadHandler interface {
	CreateThread(w http.ResponseWriter, r *http.Request)
	GetDetails(w http.ResponseWriter, r *http.Request)
	UpdateDetails(w http.ResponseWriter, r *http.Request)
	GetPosts(w http.ResponseWriter, r *http.Request)
	Vote(w http.ResponseWriter, r *http.Request)
}

type ThreadRepo interface {
	CreateThread(ctx context.Context, thread *models.Thread) (int, error)
	UpdateThreadBySlug(ctx context.Context, thread *models.Thread) error
	UpdateVote(ctx context.Context, vote *models.Vote) error
	AddVote(ctx context.Context, vote *models.Vote) error
	GetThreadBySlugOrId(ctx context.Context, slugOrId string) (*models.Thread, error)
	GetPosts(ctx context.Context, threadPosts *models.ThreadPosts) (*[]models.Post, error)
}
