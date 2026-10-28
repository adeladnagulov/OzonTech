package additem

import (
	"encoding/json"
	"errors"
	"fmt"
	"hw_2/internal/cart/domain/cart"
	"hw_2/internal/cart/http/response"
	"net/http"
	"strconv"
)

type Handler struct {
	uc CartUseCase
}

func NewHandler(uc CartUseCase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) AddItem(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.PathValue("user_id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		response.Problem(w, http.StatusInternalServerError, fmt.Errorf("failed parse user_id: %w", err))
		return
	}
	if userId <= 0 {
		response.Problem(w, http.StatusBadRequest, cart.ErrUserIdBadRequest)
		return
	}
	skuIdStr := r.PathValue("sku_id")
	skuId, err := strconv.Atoi(skuIdStr)
	if err != nil {
		response.Problem(w, http.StatusInternalServerError, fmt.Errorf("failed parse sku_id: %w", err))
		return
	}
	if skuId <= 0 {
		response.Problem(w, http.StatusBadRequest, cart.ErrSkuIdBadRequest)
		return
	}

	var req cart.AddItemReq
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.Problem(w, http.StatusInternalServerError, fmt.Errorf("failed decode addItem body: %w", err))
		return
	}

	err = h.uc.AddItem(req, r.Context(), int64(skuId), int64(userId))
	if err != nil {
		switch {
		case errors.Is(err, cart.ErrBadRequest):
			response.Problem(w, http.StatusBadRequest, err)
		case errors.Is(err, cart.ErrProductNotFound):
			response.Problem(w, http.StatusPreconditionFailed, err)
		default:
			response.Problem(w, http.StatusInternalServerError, err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
