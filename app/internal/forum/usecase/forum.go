package usecase

import (
	"context"
	"net/http"

	forumModel "github.com/forums/app/internal/forum"
	userModel "github.com/forums/app/internal/user"
	"github.com/forums/app/models"
	"github.com/forums/utils/response"
)

type usecase struct {
	forumRepo forumModel.ForumRepo
	userRepo  userModel.UserRepo
}

func NewForumUsecase(forumRepo forumModel.ForumRepo, userRepo userModel.UserRepo) forumModel.ForumUsecase {
	return &usecase{
		forumRepo: forumRepo,
		userRepo:  userRepo,
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

	user, err := u.userRepo.GetUserByName(ctx, forum.User)
	if err == nil && user == nil {
		message := models.Message{
			Message: "Can't find user with id #" + forum.User + "\n",
		}
		response := response.New(http.StatusNotFound, message)
		return response, nil
	}
	// if user != nil {
	// 	return nil, err
	// }

	forum.User = user.Nickname
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

func (u *usecase) GetUsers(ctx context.Context, forumUsers models.ForumUsers) (response.Response, error) {
	forum, err := u.forumRepo.GetForumBySlug(ctx, forumUsers.Slug)
	if err != nil {
		return nil, err
	}

	if forum == nil {
		message := models.Message{
			Message: "Can't find forum with id #" + forumUsers.Slug + "\n",
		}
		response := response.New(http.StatusNotFound, message)
		return response, nil
	}

	users, err := u.forumRepo.GetUsers(ctx, forumUsers)
	if err != nil {
		return nil, err
	}

	response := response.New(http.StatusOK, users)
	return response, nil
}

func (u *usecase) GetThreads(ctx context.Context, forumThreads models.ForumThreads) (response.Response, error) {
	forum, err := u.forumRepo.GetForumBySlug(ctx, forumThreads.Slug)
	if err != nil {
		return nil, err
	}

	if forum == nil {
		message := models.Message{
			Message: "Can't find forum with id #" + forumThreads.Slug + "\n",
		}
		response := response.New(http.StatusNotFound, message)
		return response, nil
	}

	threads, err := u.forumRepo.GetThreads(ctx, forumThreads)
	if err != nil {
		return nil, err
	}

	response := response.New(http.StatusOK, threads)
	return response, nil
}
