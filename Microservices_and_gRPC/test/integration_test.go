//go:build api

package test

import (
	"encoding/json"
	"fmt"
	"hw_2/internal/cart/domain/cart"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration(t *testing.T) {
	t.Parallel()
	userId := time.Now().UnixNano()
	baseURL := fmt.Sprintf("http://localhost:8080/user/%d", userId)

	t.Run("succes add and get cart", func(t *testing.T) {
		t.Parallel()

		validSKU := int64(135466315)
		addURL := fmt.Sprintf("%s/cart/%d", baseURL, validSKU)
		reqBody := map[string]int{"count": 2}
		resp := doHTTPRequest(t, http.MethodPost, addURL, reqBody)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		getURL := fmt.Sprintf("%s/cart", baseURL)
		resp = doHTTPRequest(t, http.MethodGet, getURL, nil)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var cartResp cart.Cart
		err := json.NewDecoder(resp.Body).Decode(&cartResp)
		require.NoError(t, err)

		assert.Len(t, cartResp.Items, 1)
		assert.Equal(t, validSKU, cartResp.Items[0].Sku)
		assert.Equal(t, uint32(2), cartResp.Items[0].Count)
		assert.Equal(t, uint32(6644), cartResp.TotalPrice)
	})

	t.Run("non existent product", func(t *testing.T) {
		t.Parallel()

		addURL := fmt.Sprintf("%s/cart/%d", baseURL, 2)
		reqBody := map[string]int{"count": 1}

		resp := doHTTPRequest(t, http.MethodPost, addURL, reqBody)
		defer resp.Body.Close()

		require.Equal(t, http.StatusPreconditionFailed, resp.StatusCode)
	})

	t.Run("error get empty cart", func(t *testing.T) {
		t.Parallel()

		emptyUserID := time.Now().UnixNano()
		getURL := fmt.Sprintf("http://localhost:8080/user/%d/cart", emptyUserID)

		resp := doHTTPRequest(t, http.MethodGet, getURL, nil)
		defer resp.Body.Close()

		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
