package delivery

import (
	"encoding/json"
	"net/http"

	threadModel "github.com/forums/app/internal/thread"
	"github.com/forums/app/models"
	"github.com/forums/utils/errors"
	"github.com/forums/utils/logger"
	"github.com/gorilla/mux"
)

type Handler struct {
	threadUsecase threadModel.ThreadUsecase
}

func NewThreadHandler(usecase threadModel.ThreadUsecase) threadModel.ThreadHandler {
	return &Handler{
		threadUsecase: usecase,
	}
}

func (h *Handler) CreateThread(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	slug := vars["slug"]
	newThread := new(models.Thread)
	err := json.NewDecoder(r.Body).Decode(&newThread)
	if err != nil {
		sendErr := errors.New(http.StatusBadRequest, err.Error())
		logger.Delivery().Error(ctx, sendErr)
		w.WriteHeader(sendErr.Code())
		return
	}
	defer r.Body.Close()
	logger.Delivery().Info(ctx, logger.Fields{"request data": *newThread})

	response, err := h.threadUsecase.CreateThread(ctx, newThread, slug)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	(*response).SendSuccess(w)
}

func (h *Handler) GetDetails(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	slugOrId := vars["slug_or_id"]
	logger.Delivery().Info(ctx, logger.Fields{"request data slug or id": slugOrId})

	response, err := h.threadUsecase.GetThread(ctx, slugOrId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	(*response).SendSuccess(w)
}

func (h *Handler) UpdateDetails(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	slugOrId := vars["slug_or_id"]
	newThread := new(models.Thread)
	err := json.NewDecoder(r.Body).Decode(&newThread)
	if err != nil {
		sendErr := errors.New(http.StatusBadRequest, err.Error())
		logger.Delivery().Error(ctx, sendErr)
		w.WriteHeader(sendErr.Code())
		return
	}
	defer r.Body.Close()
	logger.Delivery().Info(ctx, logger.Fields{"request data": *newThread})

	response, err := h.threadUsecase.UpdateThread(ctx, newThread, slugOrId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	(*response).SendSuccess(w)
}

func (h *Handler) GetPosts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	threadPosts := new(models.ThreadPosts)

	threadPosts.SlugOrId = vars["slug_or_id"]
	threadPosts.Limit = r.URL.Query().Get("limit")
	threadPosts.Since = r.URL.Query().Get("since")
	threadPosts.Sort = r.URL.Query().Get("sort")

	desc := r.URL.Query().Get("desc")
	if desc == "false" || desc == "" {
		threadPosts.Desc = false
	} else {
		threadPosts.Desc = true
	}

	logger.Delivery().Info(ctx, logger.Fields{"request data": *threadPosts})

	response, err := h.threadUsecase.GetPosts(ctx, threadPosts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	(*response).SendSuccess(w)
}

func (h *Handler) Vote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	slugOrId := vars["slug_or_id"]
	vote := new(models.Vote)
	err := json.NewDecoder(r.Body).Decode(&vote)
	if err != nil {
		sendErr := errors.New(http.StatusBadRequest, err.Error())
		logger.Delivery().Error(ctx, sendErr)
		w.WriteHeader(sendErr.Code())
		return
	}
	defer r.Body.Close()
	logger.Delivery().Info(ctx, logger.Fields{"request data": *vote})

	response, err := h.threadUsecase.AddVote(ctx, vote, slugOrId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	(*response).SendSuccess(w)
}
