package usecase

import (
	"context"
	"net/http"

	forumModel "github.com/forums/app/internal/forum"
	"github.com/forums/app/models"
	"github.com/forums/utils/response"
)

type usecase struct {
	forumRepo forumModel.UserRepo
}

func NewForumUsecase(forumRepo forumModel.UserRepo) forumModel.ForumUsecase {
	return &usecase{
		forumRepo: forumRepo,
	}
}

func (u *usecase) CreateForum(ctx context.Context, forum models.Forum) (response.Response, error) {

	forumDb, err := u.forumRepo.GetForumBySlug(ctx, forum.Slug)
	if err != nil {
		return nil, err
	}
	if forumDb != nil {
		response := response.New(http.StatusConflict, forumDb)
		return response, nil
	}

	_, err = u.forumRepo.CreateForum(ctx, forum)
	if err != nil {
		return nil, err
	}

	response := response.New(http.StatusCreated, forum)
	return response, nil
}

func (u *usecase) GetForumBySlug(ctx context.Context, slug string) (response.Response, error) {

	forum, err := u.forumRepo.GetForumBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if forum == nil {
		message := models.Message{
			Message: "Can't find forum with id #" + slug + "\n",
		}
		response := response.New(http.StatusNotFound, message)
		return response, nil
	}

	response := response.New(http.StatusOK, forum)
	return response, nil
}
