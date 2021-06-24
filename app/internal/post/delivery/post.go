package delivery

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	forumModel "github.com/forums/app/internal/forum"
	postModel "github.com/forums/app/internal/post"
	threadModel "github.com/forums/app/internal/thread"
	userModel "github.com/forums/app/internal/user"
	"github.com/forums/app/models"
	"github.com/forums/utils/errors"
	"github.com/forums/utils/logger"
	"github.com/forums/utils/response"
	"github.com/gorilla/mux"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
)

type Handler struct {
	postRepo   postModel.PostRepo
	userRepo   userModel.UserRepo
	threadRepo threadModel.ThreadRepo
	forumRepo  forumModel.ForumRepo
}

func NewPostHandler(postRepo postModel.PostRepo, userRepo userModel.UserRepo,
	threadRepo threadModel.ThreadRepo, forumRepo forumModel.ForumRepo) postModel.PostHandler {
	return &Handler{
		postRepo:   postRepo,
		userRepo:   userRepo,
		threadRepo: threadRepo,
		forumRepo:  forumRepo,
	}
}

func (h *Handler) CreatePosts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	logger.Delivery().AddFuncName("CreatePosts").InlineDebug(ctx, "request")

	posts := make([]models.Post, 0)
	// с этим работает при массивах и пустых тоже
	err := json.NewDecoder(r.Body).Decode(&posts)
	if err != nil {
		sendErr := errors.New(http.StatusBadRequest, err.Error())
		logger.Delivery().Error(ctx, sendErr)
		w.WriteHeader(sendErr.Code())
		return
	}
	defer r.Body.Close()
	slug := vars["slug_or_id"]
	logger.Delivery().Info(ctx, logger.Fields{"request data": posts, "slug_or_id": slug})

	timeNow := time.Now()

	thread, err := h.threadRepo.GetThreadBySlugOrId(ctx, slug)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if thread == nil {
		message := models.Message{
			Message: "Can't find thread with id #" + slug + "\n",
		}
		response.New(http.StatusNotFound, message).SendSuccess(w)
		return
	}

	if len(posts) == 0 {
		response.New(http.StatusCreated, posts).SendSuccess(w)
		return
	}

	logger.Usecase().Debug(ctx, logger.Fields{"forum slug": thread.Forum})
	for i := range posts {
		posts[i].Thread = thread.Id
		posts[i].Forum = thread.Forum
		posts[i].Created = timeNow
	}

	postsDB, err := h.postRepo.CreatePosts(ctx, &posts)
	if err != nil {
		logger.Usecase().AddFuncName("CreatePosts").Info(ctx, logger.Fields{"Error": err})
		if pqErr, ok := err.(pgx.PgError); ok {
			logger.Usecase().AddFuncName("CreatePosts").Info(ctx, logger.Fields{"Error Code": pqErr.Code})
			switch pqErr.Code {
			case pgerrcode.ForeignKeyViolation: // проблемы с сохранением user
				message := models.Message{
					Message: "Can't find user\n",
				}
				response.New(http.StatusNotFound, message).SendSuccess(w)
				return

			case "12345":
				{
					message := models.Message{
						Message: "Parent not found\n",
					}
					response.New(http.StatusConflict, message).SendSuccess(w)
					return
				}

			default:
				logger.Usecase().AddFuncName("CreatePosts").Error(ctx, err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			logger.Usecase().AddFuncName("CreatePosts").Info(ctx, logger.Fields{"Error": err})
		}
	}

	response.New(http.StatusCreated, postsDB).SendSuccess(w)
}

func (h *Handler) GetDetails(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	related := new(models.RequestPost)
	related.Related = r.URL.Query().Get("related")

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendErr := errors.New(http.StatusBadRequest, err.Error())
		logger.Delivery().Error(ctx, sendErr)
		w.WriteHeader(sendErr.Code())
		return
	}
	related.Id = id
	logger.Delivery().Info(ctx, logger.Fields{"request data": *related})

	// TODO: подумать надо оптимизацией, получениявсех данных одним запросом
	post, err := h.postRepo.GetPost(ctx, related.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if post == nil {
		message := models.Message{
			Message: "Can't find post with id #" + strconv.Itoa(related.Id) + "\n",
		}
		response.New(http.StatusNotFound, message).SendSuccess(w)
		return
	}

	infoPost := models.InfoPost{
		Post:   post,
		User:   nil,
		Forum:  nil,
		Thread: nil,
	}

	if strings.Contains(related.Related, "user") {
		user, err := h.userRepo.GetUserByName(ctx, post.Author)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		infoPost.User = user
	}

	if strings.Contains(related.Related, "forum") {
		forum, err := h.forumRepo.GetForumBySlug(ctx, post.Forum)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		infoPost.Forum = forum
	}

	if strings.Contains(related.Related, "thread") {
		thread, err := h.threadRepo.GetThreadBySlugOrId(ctx, strconv.Itoa(post.Thread))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		infoPost.Thread = thread
	}

	response.New(http.StatusOK, infoPost).SendSuccess(w)
}

func (h *Handler) UpdateDetails(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	message := new(models.MessagePostRequest)
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		sendErr := errors.New(http.StatusBadRequest, err.Error())
		logger.Delivery().Error(ctx, sendErr)
		w.WriteHeader(sendErr.Code())
		return
	}
	defer r.Body.Close()
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendErr := errors.New(http.StatusBadRequest, err.Error())
		logger.Delivery().Error(ctx, sendErr)
		w.WriteHeader(sendErr.Code())
		return
	}
	message.Id = id
	logger.Delivery().Info(ctx, logger.Fields{"request data": *message})

	post, err := h.postRepo.GetPost(ctx, message.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if post == nil {
		message := models.Message{
			Message: "Can't find post with id #" + strconv.Itoa(message.Id) + "\n",
		}
		response.New(http.StatusNotFound, message).SendSuccess(w)
		return
	}

	if post.Message == message.Message || message.Message == "" {
		response.New(http.StatusOK, post).SendSuccess(w)
		return
	}

	post.Message = message.Message
	post.IsEdited = true

	err = h.postRepo.UpdateMessage(ctx, message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.New(http.StatusOK, post).SendSuccess(w)
}
