package cart

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) AddItem(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.PathValue("user_id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Errorf("failed parse user_id: %w", err))
		return
	}
	if userId <= 0 {
		writeError(w, http.StatusBadRequest, ErrUserIdBadRequest)
		return
	}
	skuIdStr := r.PathValue("sku_id")
	skuId, err := strconv.Atoi(skuIdStr)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Errorf("failed parse sku_id: %w", err))
		return
	}
	if skuId <= 0 {
		writeError(w, http.StatusBadRequest, ErrSkuIdBadRequest)
		return
	}

	var req AddItemReq
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Errorf("failed decode addItem body: %w", err))
		return
	}

	err = h.service.AddItem(req, r.Context(), int64(skuId), int64(userId))
	if err != nil {
		switch {
		case errors.Is(err, ErrBadRequest):
			writeError(w, http.StatusBadRequest, err)
		case errors.Is(err, ErrProductNotFound):
			writeError(w, http.StatusPreconditionFailed, err)
		default:
			writeError(w, http.StatusInternalServerError, err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Service) GetCart(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func writeError(w http.ResponseWriter, code int, desc error) {
	var problem struct {
		Code int    `json:"code"`
		Desc string `json:"description"`
	}
	problem.Code = code
	problem.Desc = desc.Error()
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(problem)
}
