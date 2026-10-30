package get

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
	us CartUseCase
}

func NewHandler(us CartUseCase) *Handler {
	return &Handler{
		us: us,
	}
}

func (h *Handler) GetCart(w http.ResponseWriter, r *http.Request) {
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

	c, err := h.us.GetCart(r.Context(), int64(userId))
	if err != nil {
		if errors.Is(err, cart.ErrProductNotFound) {
			response.Problem(w, http.StatusNotFound, fmt.Errorf("failed get cart by user_id = %d: %w", userId, err))
			return
		} else if errors.Is(err, cart.ErrProductNotFound) {
			response.Problem(w, http.StatusPreconditionFailed, err)
			return
		}
		response.Problem(w, http.StatusInternalServerError, fmt.Errorf("failed get cart: %w", err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(c)
}
