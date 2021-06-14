package delivery

import (
	"encoding/json"
	"net/http"
	"strconv"

	postModel "github.com/forums/app/internal/post"
	"github.com/forums/app/models"
	"github.com/forums/utils/errors"
	"github.com/forums/utils/logger"
	"github.com/gorilla/mux"
)

type Handler struct {
	postUsecase postModel.PostUsecase
}

func NewPostHandler(postUsecase postModel.PostUsecase) postModel.PostHandler {
	return &Handler{
		postUsecase: postUsecase,
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

	response, err := h.postUsecase.CreatePosts(ctx, posts, slug)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w)
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

	response, err := h.postUsecase.GetDetails(ctx, *related)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w)
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

	response, err := h.postUsecase.UpdateMessage(ctx, *message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w)
}
