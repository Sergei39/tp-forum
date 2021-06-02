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

func (u *usecase) GetThread(ctx context.Context, slug_or_id string) (response.Response, error) {
	// TODO: понять что все таки может прийти в запросе slug или id
	thread, err := u.threadRepo.GetThreadBySlug(ctx, slug_or_id)
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

	response := response.New(http.StatusCreated, thread)
	return response, nil
}

func (u *usecase) UpdateThread(ctx context.Context, thread models.Thread, slugOrId string) (response.Response, error) {
	// TODO: понять что все таки может прийти в запросе slug или id
	threadOld, err := u.threadRepo.GetThreadBySlug(ctx, slugOrId)
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

	threadOld.Title = thread.Title
	threadOld.Message = thread.Message

	err = u.threadRepo.UpdateThreadBySlug(ctx, *threadOld)
	if err != nil {
		return nil, err
	}

	response := response.New(http.StatusCreated, threadOld)
	return response, nil
}

func (u *usecase) AddVote(ctx context.Context, vote models.Vote, slugOrId string) (response.Response, error) {
	// TODO: понять что все таки может прийти в запросе slug или id
	// TODO: подумать как это сделать меньшим кол-вом запросов
	thread, err := u.threadRepo.GetThreadBySlug(ctx, slugOrId)
	if err != nil {
		return nil, err
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

	thread, err = u.threadRepo.GetThreadBySlug(ctx, slugOrId)
	if err != nil {
		return nil, err
	}
	response := response.New(http.StatusCreated, thread)
	return response, nil
}
