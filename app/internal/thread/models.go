package thread

import (
	"context"

	"github.com/forums/app/models"
	"github.com/forums/utils/response"

	"github.com/labstack/echo/v4"
)

type ThreadHandler interface {
	CreateThread(c echo.Context) error
	GetDetails(c echo.Context) error
	UpdateDetails(c echo.Context) error
	GetPosts(c echo.Context) error
	Vote(c echo.Context) error
}

type ThreadUsecase interface {
	CreateThread(ctx context.Context, thread models.Thread, slug string) (response.Response, error)
	GetThread(ctx context.Context, slug_or_id string) (response.Response, error)
	UpdateThread(ctx context.Context, thread models.Thread, slugOrId string) (response.Response, error)
	AddVote(ctx context.Context, vote models.Vote, slugOrId string) (response.Response, error)
}

type ThreadRepo interface {
	CreateThread(ctx context.Context, thread models.Thread) (int, error)
	GetThreadBySlug(ctx context.Context, slug string) (*models.Thread, error)
	UpdateThreadBySlug(ctx context.Context, thread models.Thread) error
	UpdateVote(ctx context.Context, vote models.Vote) error
	CheckVote(ctx context.Context, vote models.Vote) (int, bool, error)
	AddVote(ctx context.Context, vote models.Vote) error
}
