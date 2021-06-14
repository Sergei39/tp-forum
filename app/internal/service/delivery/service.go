package delivery

import (
	"net/http"

	serviceModel "github.com/forums/app/internal/service"
)

type Handler struct {
	serviceUsecase serviceModel.ServiceUsecase
}

func NewServiceHandler(serviceUsecase serviceModel.ServiceUsecase) serviceModel.ServiceHandler {
	return &Handler{
		serviceUsecase: serviceUsecase,
	}
}

func (h *Handler) ClearDb(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err := h.serviceUsecase.ClearDb(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) StatusDb(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	response, err := h.serviceUsecase.StatusDb(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.SendSuccess(w)
}
