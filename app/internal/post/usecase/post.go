package usecase

import (
	"context"
	"net/http"
	"strconv"
	"time"

	forumModel "github.com/forums/app/internal/forum"
	postModel "github.com/forums/app/internal/post"
	threadModel "github.com/forums/app/internal/thread"
	userModel "github.com/forums/app/internal/user"
	"github.com/forums/app/models"
	"github.com/forums/utils/logger"
	"github.com/forums/utils/response"
)

type usecase struct {
	postRepo   postModel.PostRepo
	userRepo   userModel.UserRepo
	threadRepo threadModel.ThreadRepo
	forumRepo  forumModel.ForumRepo
}

func NewPostUsecase(postRepo postModel.PostRepo, userRepo userModel.UserRepo,
	threadRepo threadModel.ThreadRepo, forumRepo forumModel.ForumRepo) postModel.PostUsecase {
	return &usecase{
		postRepo:   postRepo,
		userRepo:   userRepo,
		threadRepo: threadRepo,
		forumRepo:  forumRepo,
	}
}

func (u *usecase) GetDetails(ctx context.Context, request models.RequestPost) (response.Response, error) {
	// TODO: доделать содержание ответа в зависимости от параметров запроса
	// TODO: вернуться к этому запросу и дописать thread
	// TODO: подумать надо оптимизацией, получениявсех данных одним запросом
	post, err := u.postRepo.GetPost(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	if post == nil {
		message := models.Message{
			Message: "Can't find post with id #" + strconv.Itoa(request.Id) + "\n",
		}
		response := response.New(http.StatusNotFound, message)
		return response, nil
	}

	user, err := u.userRepo.GetUserByName(ctx, post.Author)
	if err != nil {
		return nil, err
	}

	forum, err := u.forumRepo.GetForumBySlug(ctx, post.Forum)
	if err != nil {
		return nil, err
	}

	infoPost := models.InfoPost{
		Post:  *post,
		User:  *user,
		Forum: *forum,
	}

	response := response.New(http.StatusOK, infoPost)
	return response, nil
}

func (u *usecase) UpdateMessage(ctx context.Context, request models.MessagePostRequest) (response.Response, error) {
	post, err := u.postRepo.GetPost(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	if post == nil {
		message := models.Message{
			Message: "Can't find post with id #" + strconv.Itoa(request.Id) + "\n",
		}
		response := response.New(http.StatusNotFound, message)
		return response, nil
	}

	if post.Message == request.Message {
		response := response.New(http.StatusOK, post)
		return response, nil
	}

	post.Message = request.Message
	post.IsEdited = true

	err = u.postRepo.UpdateMessage(ctx, request)
	if err != nil {
		return nil, err
	}

	response := response.New(http.StatusOK, post)
	return response, nil
}

func (u *usecase) CreatePosts(ctx context.Context, posts []models.Post, slugOrId string) (response.Response, error) {
	timeNow := time.Now()

	if len(posts) == 0 {
		response := response.New(http.StatusCreated, posts)
		return response, nil
	}

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

	logger.Usecase().Debug(ctx, logger.Fields{"forum slug": thread.Forum})
	for i := range posts {
		posts[i].Thread = thread.Id
		posts[i].Forum = thread.Forum
		posts[i].Created = timeNow
		nest, err := u.createTreeArray(ctx, int(posts[i].Parent))
		if err != nil {
			return nil, err
		}

		id, err := u.postRepo.CreatePost(ctx, posts[i], nest)
		if err != nil {
			return nil, err
		}

		posts[i].Id = int64(id)
	}

	// TODO: сделать проверку на то что parent валидный

	response := response.New(http.StatusCreated, posts)
	return response, nil
}

func (u *usecase) createTreeArray(ctx context.Context, id int) ([]int64, error) {
	nest, err := u.postRepo.GetPostAndChildLastArr(ctx, id)
	if err != nil {
		return nil, err
	}

	if len(nest.Last) != 0 {
		nest.Last[len(nest.Last)-1]++
		return nest.Last, nil
	}

	if len(nest.Parent) != 0 {
		tecNest := nest.Parent
		tecNest = append(tecNest, 1)
		return tecNest, nil
	}

	tecNest := make([]int64, 0)
	tecNest = append(tecNest, 1)
	return tecNest, nil
}
