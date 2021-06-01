package usecase

import (
	"context"
	"net/http"

	forumModel "github.com/forums/app/internal/forum"
	threadModel "github.com/forums/app/internal/thread"
	userModel "github.com/forums/app/internal/user"
	"github.com/forums/app/models"
	"github.com/forums/utils/response"
)

type usecase struct {
	threadRepo threadModel.ThreadRepo
	userRepo   userModel.UserRepo
	forumRepo  forumModel.ForumRepo
}

func NewThreadUsecase(threadRepo threadModel.ThreadRepo, userRepo userModel.UserRepo,
	forumRepo forumModel.ForumRepo) threadModel.ThreadUsecase {
	return &usecase{
		threadRepo: threadRepo,
		userRepo:   userRepo,
		forumRepo:  forumRepo,
	}
}

func (u *usecase) CreateThread(ctx context.Context, thread models.Thread, slug string) (response.Response, error) {
	user, err := u.userRepo.GetUserByName(ctx, thread.Author)
	if err != nil {
		return nil, err
	}
	forum, err := u.forumRepo.GetForumBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	if user == nil {
		message := models.Message{
			Message: "Can't find user with id #" + thread.Author + "\n",
		}
		response := response.New(http.StatusNotFound, message)
		return response, nil
	}
	if forum == nil {
		message := models.Message{
			Message: "Can't find forum with id #" + slug + "\n",
		}
		response := response.New(http.StatusNotFound, message)
		return response, nil
	}

	// TODO: проверка на существования ветки уже
	id, err := u.threadRepo.CreateThread(ctx, thread)
	if err != nil {
		return nil, err
	}

	thread.Id = id
	response := response.New(http.StatusCreated, thread)
	return response, nil
}
