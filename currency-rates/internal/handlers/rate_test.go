package handlers_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/RinaDish/currency-rates/internal/handlers"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type successRateClient struct{}

func (c successRateClient) GetDollarRate(ctx context.Context) (float64, error) {
	return 10.0, nil
}

type failedRateClient struct{}

func (c failedRateClient) GetDollarRate(ctx context.Context) (float64, error) {
	return 0.0, errors.New("banks not available")
}

func TestGetCurrentRate(main *testing.T) {
	l := zap.NewNop()
	main.Run("succesfully returned rates", func(t *testing.T) {
		expectedRate := float64(10.0)
		s := successRateClient{}
		h := handlers.NewRateHandler(l.Sugar(), s)

		req := httptest.NewRequest(http.MethodGet, "/rates", nil)
		w := httptest.NewRecorder()

		h.GetCurrentRate(w, req)
		res := w.Result()
		defer res.Body.Close()
		data, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		rate, err := strconv.ParseFloat(string(data), 64)
		require.NoError(t, err)
		require.Equal(t, expectedRate, rate)

		require.Equal(t, w.Result().StatusCode, http.StatusOK)
	})

	main.Run("failed returned rates", func(t *testing.T) {
		f := failedRateClient{}
		h := handlers.NewRateHandler(l.Sugar(), f)

		req := httptest.NewRequest(http.MethodGet, "/rates", nil)
		w := httptest.NewRecorder()

		h.GetCurrentRate(w, req)
		res := w.Result()
		defer res.Body.Close()
		data, err := io.ReadAll(res.Body)

		require.NoError(t, err)
		require.Empty(t, data)

		require.Equal(t, w.Result().StatusCode, http.StatusBadRequest)
	})
}
