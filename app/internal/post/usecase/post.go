package usecase

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	forumModel "github.com/forums/app/internal/forum"
	postModel "github.com/forums/app/internal/post"
	threadModel "github.com/forums/app/internal/thread"
	userModel "github.com/forums/app/internal/user"
	"github.com/forums/app/models"
	"github.com/forums/utils/logger"
	"github.com/forums/utils/response"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
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

func (u *usecase) GetDetails(ctx context.Context, request *models.RequestPost) (*response.Response, error) {
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
		return &response, nil
	}

	infoPost := models.InfoPost{
		Post:   post,
		User:   nil,
		Forum:  nil,
		Thread: nil,
	}

	if strings.Contains(request.Related, "user") {
		user, err := u.userRepo.GetUserByName(ctx, post.Author)
		if err != nil {
			return nil, err
		}
		infoPost.User = user
	}

	if strings.Contains(request.Related, "forum") {
		forum, err := u.forumRepo.GetForumBySlug(ctx, post.Forum)
		if err != nil {
			return nil, err
		}
		infoPost.Forum = forum
	}

	if strings.Contains(request.Related, "thread") {
		thread, err := u.threadRepo.GetThreadBySlugOrId(ctx, strconv.Itoa(post.Thread))
		if err != nil {
			return nil, err
		}
		infoPost.Thread = thread
	}

	response := response.New(http.StatusOK, infoPost)
	return &response, nil
}

func (u *usecase) UpdateMessage(ctx context.Context, request *models.MessagePostRequest) (*response.Response, error) {
	post, err := u.postRepo.GetPost(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	if post == nil {
		message := models.Message{
			Message: "Can't find post with id #" + strconv.Itoa(request.Id) + "\n",
		}
		response := response.New(http.StatusNotFound, message)
		return &response, nil
	}

	if post.Message == request.Message || request.Message == "" {
		response := response.New(http.StatusOK, post)
		return &response, nil
	}

	post.Message = request.Message
	post.IsEdited = true

	err = u.postRepo.UpdateMessage(ctx, request)
	if err != nil {
		return nil, err
	}

	response := response.New(http.StatusOK, post)
	return &response, nil
}

func (u *usecase) CreatePosts(ctx context.Context, posts *[]models.Post, slugOrId string) (*response.Response, error) {
	timeNow := time.Now()

	thread, err := u.threadRepo.GetThreadBySlugOrId(ctx, slugOrId)
	if err != nil {
		return nil, err
	}
	if thread == nil {
		message := models.Message{
			Message: "Can't find thread with id #" + slugOrId + "\n",
		}
		response := response.New(http.StatusNotFound, message)
		return &response, nil
	}

	if len(*posts) == 0 {
		response := response.New(http.StatusCreated, posts)
		return &response, nil
	}

	logger.Usecase().Debug(ctx, logger.Fields{"forum slug": thread.Forum})
	for i := range *posts {
		(*posts)[i].Thread = thread.Id
		(*posts)[i].Forum = thread.Forum
		(*posts)[i].Created = timeNow
	}

	postsDB, err := u.postRepo.CreatePosts(ctx, posts)
	if err != nil {
		logger.Usecase().AddFuncName("CreatePosts").Info(ctx, logger.Fields{"Error": err})
		if pqErr, ok := err.(pgx.PgError); ok {
			logger.Usecase().AddFuncName("CreatePosts").Info(ctx, logger.Fields{"Error Code": pqErr.Code})
			switch pqErr.Code {
			case pgerrcode.ForeignKeyViolation: // проблемы с сохранением user
				message := models.Message{
					Message: "Can't find user\n",
				}
				response := response.New(http.StatusNotFound, message)
				return &response, nil

			case "12345":
				{
					message := models.Message{
						Message: "Parent not found\n",
					}
					response := response.New(http.StatusConflict, message)
					return &response, nil
				}

			default:
				logger.Usecase().AddFuncName("CreatePosts").Error(ctx, err)
				return nil, err
			}
		} else {
			logger.Usecase().AddFuncName("CreatePosts").Info(ctx, logger.Fields{"Error": err})
		}
	}

	// if err = u.postRepo.CreateForumsUsers(ctx, posts); err != nil {
	// 	return nil, err
	// }

	response := response.New(http.StatusCreated, postsDB)
	return &response, nil
}
