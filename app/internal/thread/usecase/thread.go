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

	oldThread, err := u.threadRepo.GetThreadBySlugOrId(ctx, thread.Slug)
	if err != nil {
		return nil, err
	}
	if thread.Slug != "" && oldThread != nil {
		response := response.New(http.StatusConflict, oldThread)
		return response, nil
	}

	thread.Forum = forum.Slug
	id, err := u.threadRepo.CreateThread(ctx, thread)
	if err != nil {
		return nil, err
	}

	thread.Id = id
	response := response.New(http.StatusCreated, thread)
	return response, nil
}

func (u *usecase) GetThread(ctx context.Context, slug_or_id string) (response.Response, error) {
	thread, err := u.threadRepo.GetThreadBySlugOrId(ctx, slug_or_id)
	if err != nil {
		return nil, err
	}

	if thread == nil {
		message := models.Message{
			Message: "Can't find thread with id #" + slug_or_id + "\n",
		}
		response := response.New(http.StatusNotFound, message)
		return response, nil
	}

	response := response.New(http.StatusOK, thread)
	return response, nil
}

func (u *usecase) fixData(newThread, oldThread models.Thread) models.Thread {
	if newThread.Message == "" {
		newThread.Message = oldThread.Message
	}

	if newThread.Title == "" {
		newThread.Title = oldThread.Title
	}

	return newThread
}

func (u *usecase) UpdateThread(ctx context.Context, thread models.Thread, slugOrId string) (response.Response, error) {
	threadOld, err := u.threadRepo.GetThreadBySlugOrId(ctx, slugOrId)
	if err != nil {
		return nil, err
	}

	if threadOld == nil {
		message := models.Message{
			Message: "Can't find thread with id #" + slugOrId + "\n",
		}
		response := response.New(http.StatusNotFound, message)
		return response, nil
	}

	if thread.Title != "" {
		threadOld.Title = thread.Title
	}

	if thread.Message != "" {
		threadOld.Message = thread.Message
	}

	err = u.threadRepo.UpdateThreadBySlug(ctx, *threadOld)
	if err != nil {
		return nil, err
	}

	response := response.New(http.StatusOK, threadOld)
	return response, nil
}

func (u *usecase) AddVote(ctx context.Context, vote models.Vote, slugOrId string) (response.Response, error) {
	// TODO: подумать как это сделать меньшим кол-вом запросов
	// TODO: убрать проверку user и thread на бд
	thread, err := u.threadRepo.GetThreadBySlugOrId(ctx, slugOrId)
	if err != nil {
		return nil, err
	}
	if thread == nil {
		message := models.Message{
			Message: "Can't find thread with id #" + slugOrId + "\n",
		}
		response := response.New(http.StatusNotFound, message)
		return response, nil
	}

	user, err := u.userRepo.GetUserByName(ctx, vote.User)
	if err != nil {
		return nil, err
	}
	if user == nil {
		message := models.Message{
			Message: "Can't find user with id #" + vote.User + "\n",
		}
		response := response.New(http.StatusNotFound, message)
		return response, nil
	}

	vote.Thread = thread.Id
	id, ok, err := u.threadRepo.CheckVote(ctx, vote)
	if err != nil {
		return nil, err
	}
	if !ok {
		err = u.threadRepo.AddVote(ctx, vote)
		if err != nil {
			return nil, err
		}
	}

	if ok {
		vote.Id = id
		err = u.threadRepo.UpdateVote(ctx, vote)
		if err != nil {
			return nil, err
		}
	}

	thread, err = u.threadRepo.GetThreadBySlugOrId(ctx, slugOrId)
	if err != nil {
		return nil, err
	}
	response := response.New(http.StatusOK, thread)
	return response, nil
}

func (u *usecase) GetPosts(ctx context.Context, threadPosts models.ThreadPosts) (response.Response, error) {
	thread, err := u.threadRepo.GetThreadBySlugOrId(ctx, threadPosts.SlugOrId)
	if err != nil {
		return nil, err
	}

	if thread == nil {
		message := models.Message{
			Message: "Can't find forum with id #" + threadPosts.SlugOrId + "\n",
		}
		response := response.New(http.StatusNotFound, message)
		return response, nil
	}

	threadPosts.ThreadId = thread.Id
	posts, err := u.threadRepo.GetPosts(ctx, threadPosts)
	if err != nil {
		return nil, err
	}

	response := response.New(http.StatusOK, posts)
	return response, nil
}
