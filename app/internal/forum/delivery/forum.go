package delivery

import (
	"encoding/json"
	"net/http"
	"strconv"

	forumModel "github.com/forums/app/internal/forum"
	"github.com/forums/app/models"
	"github.com/forums/utils/errors"
	"github.com/forums/utils/logger"
	"github.com/gorilla/mux"
)

type Handler struct {
	forumUsecase forumModel.ForumUsecase
}

func NewForumHandler(usecase forumModel.ForumUsecase) forumModel.ForumHandler {
	return &Handler{
		forumUsecase: usecase,
	}
}

func (h *Handler) CreateForum(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	newForum := new(models.Forum)
	err := json.NewDecoder(r.Body).Decode(&newForum)
	if err != nil {
		sendErr := errors.New(http.StatusBadRequest, err.Error())
		logger.Delivery().Error(ctx, sendErr)
		w.WriteHeader(sendErr.Code())
		return
	}
	defer r.Body.Close()
	logger.Delivery().Info(ctx, logger.Fields{"request data": *newForum})

	response, err := h.forumUsecase.CreateForum(ctx, *newForum)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Delivery().Debug(ctx, logger.Fields{"response": response})
	response.SendSuccess(w)
}

func (h *Handler) GetDetails(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	slug := vars["slug"]
	logger.Delivery().Info(ctx, logger.Fields{"request data": slug})

	response, err := h.forumUsecase.GetForumBySlug(ctx, slug)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w)
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

	response, err := h.forumUsecase.GetUsers(ctx, *forumUsers)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w)
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

	response, err := h.forumUsecase.GetThreads(ctx, *forumThreads)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w)
}
