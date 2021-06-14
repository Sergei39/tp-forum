package thread

import (
	"context"
	"net/http"

	"github.com/forums/app/models"
	"github.com/forums/utils/response"
)

type ThreadHandler interface {
	CreateThread(w http.ResponseWriter, r *http.Request)
	GetDetails(w http.ResponseWriter, r *http.Request)
	UpdateDetails(w http.ResponseWriter, r *http.Request)
	GetPosts(w http.ResponseWriter, r *http.Request)
	Vote(w http.ResponseWriter, r *http.Request)
}

type ThreadUsecase interface {
	CreateThread(ctx context.Context, thread models.Thread, slug string) (response.Response, error)
	GetThread(ctx context.Context, slug_or_id string) (response.Response, error)
	UpdateThread(ctx context.Context, thread models.Thread, slugOrId string) (response.Response, error)
	AddVote(ctx context.Context, vote models.Vote, slugOrId string) (response.Response, error)
	GetPosts(ctx context.Context, threadPosts models.ThreadPosts) (response.Response, error)
}

type ThreadRepo interface {
	CreateThread(ctx context.Context, thread models.Thread) (int, error)
	UpdateThreadBySlug(ctx context.Context, thread models.Thread) error
	UpdateVote(ctx context.Context, vote models.Vote) error
	AddVote(ctx context.Context, vote models.Vote) error
	GetThreadBySlugOrId(ctx context.Context, slugOrId string) (*models.Thread, error)
	GetPosts(ctx context.Context, threadPosts models.ThreadPosts) ([]models.Post, error)
}
