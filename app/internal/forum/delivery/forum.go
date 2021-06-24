package delivery

import (
	"encoding/json"
	"net/http"
	"strconv"

	forumModel "github.com/forums/app/internal/forum"
	userModel "github.com/forums/app/internal/user"
	"github.com/forums/app/models"
	"github.com/forums/utils/errors"
	"github.com/forums/utils/logger"
	"github.com/forums/utils/response"
	"github.com/gorilla/mux"
)

type Handler struct {
	forumRepo forumModel.ForumRepo
	userRepo  userModel.UserRepo
}

func NewForumHandler(forumRepo forumModel.ForumRepo, userRepo userModel.UserRepo) forumModel.ForumHandler {
	return &Handler{
		forumRepo: forumRepo,
		userRepo:  userRepo,
	}
}

func (h *Handler) CreateForum(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	newForum := new(models.Forum)
	err := json.NewDecoder(r.Body).Decode(&newForum)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	logger.Delivery().Info(ctx, logger.Fields{"request data": *newForum})

	forumDb, err := h.forumRepo.GetForumBySlug(ctx, newForum.Slug)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if forumDb != nil {
		response.New(http.StatusConflict, forumDb).SendSuccess(w)
		return
	}

	user, err := h.userRepo.GetUserByName(ctx, newForum.User)
	if err == nil && user == nil {
		message := models.Message{
			Message: "Can't find user with id #" + newForum.User + "\n",
		}
		response.New(http.StatusNotFound, message).SendSuccess(w)
		return
	}

	newForum.User = user.Nickname
	_, err = h.forumRepo.CreateForum(ctx, newForum)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.New(http.StatusCreated, newForum).SendSuccess(w)
}

func (h *Handler) GetDetails(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	slug := vars["slug"]
	logger.Delivery().Info(ctx, logger.Fields{"request data": slug})

	forum, err := h.forumRepo.GetForumBySlug(ctx, slug)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if forum == nil {
		message := models.Message{
			Message: "Can't find forum with id #" + slug + "\n",
		}
		response.New(http.StatusNotFound, message).SendSuccess(w)
		return
	}

	response.New(http.StatusOK, forum).SendSuccess(w)
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	forumUsers := new(models.ForumUsers)

	vars := mux.Vars(r)
	forumUsers.Slug = vars["slug"]
	forumUsers.Limit = r.URL.Query().Get("limit")

	forumUsers.Since = r.URL.Query().Get("since")
	desc := r.URL.Query().Get("desc")
	if desc == "false" || desc == "" {
		forumUsers.Desc = false
	} else {
		forumUsers.Desc = true
	}

	logger.Delivery().Info(ctx, logger.Fields{"request data": *forumUsers})

	forum, err := h.forumRepo.GetForumBySlug(ctx, forumUsers.Slug)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if forum == nil {
		message := models.Message{
			Message: "Can't find forum with id #" + forumUsers.Slug + "\n",
		}
		response.New(http.StatusNotFound, message).SendSuccess(w)
		return
	}

	users, err := h.forumRepo.GetUsers(ctx, forumUsers)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.New(http.StatusOK, users).SendSuccess(w)
}

func (h *Handler) GetThreads(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	forumThreads := new(models.ForumThreads)

	vars := mux.Vars(r)
	forumThreads.Slug = vars["slug"]
	limit := r.URL.Query().Get("limit")
	if limit != "" {
		limitConv, err := strconv.Atoi(limit)
		if err != nil {
			sendErr := errors.New(http.StatusBadRequest, "convert request data - limit")
			logger.Delivery().Error(ctx, sendErr)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		forumThreads.Limit = limitConv
	}

	forumThreads.Since = r.URL.Query().Get("since")
	desc := r.URL.Query().Get("desc")
	if desc == "false" || desc == "" {
		forumThreads.Desc = false
	} else {
		forumThreads.Desc = true
	}

	logger.Delivery().Info(ctx, logger.Fields{"request data": *forumThreads})

	forum, err := h.forumRepo.GetForumBySlug(ctx, forumThreads.Slug)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if forum == nil {
		message := models.Message{
			Message: "Can't find forum with id #" + forumThreads.Slug + "\n",
		}
		response.New(http.StatusNotFound, message).SendSuccess(w)
		return
	}

	threads, err := h.forumRepo.GetThreads(ctx, forumThreads)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.New(http.StatusOK, threads).SendSuccess(w)
}
