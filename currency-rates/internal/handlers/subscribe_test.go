package handlers_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/RinaDish/currency-rates/internal/handlers"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type successDb struct{}

func (d successDb) SetEmail(ctx context.Context, email string) error {
	return nil
}


type failDb struct{}

func (d failDb) SetEmail(ctx context.Context, email string) error {
	return errors.New("email exist")
}

func TestCreateSubscription(main *testing.T) {
	l := zap.NewNop()
	main.Run("succesfully", func(t *testing.T) {
		d := successDb{}
		h := handlers.NewSubscribeHandler(l.Sugar(), d)

		form := url.Values{}
		form.Add("email", "test@test.com")

		req := httptest.NewRequest(http.MethodPost, "/subscribe", strings.NewReader(form.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		w := httptest.NewRecorder()

		h.CreateSubscription(w, req)
		res := w.Result()
		defer res.Body.Close()
		data, err := io.ReadAll(res.Body)

		require.NoError(t, err)
		require.Empty(t, data)

		require.Equal(t, w.Result().StatusCode, http.StatusOK)
	})
	main.Run("failure invalid email", func(t *testing.T) {
		d := successDb{}
		h := handlers.NewSubscribeHandler(l.Sugar(), d)

		form := url.Values{}
		form.Add("email", "test")

		req := httptest.NewRequest(http.MethodPost, "/subscribe", strings.NewReader(form.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		w := httptest.NewRecorder()

		h.CreateSubscription(w, req)
		res := w.Result()
		defer res.Body.Close()
		data, err := io.ReadAll(res.Body)

		require.NoError(t, err)
		require.Equal(t, "Invalid email\n", string(data))

		require.Equal(t, w.Result().StatusCode, http.StatusConflict)
	})
	main.Run("failure db set email", func(t *testing.T) {
		d := failDb{}
		h := handlers.NewSubscribeHandler(l.Sugar(), d)

		form := url.Values{}
		form.Add("email", "test@gmail.com")

		req := httptest.NewRequest(http.MethodPost, "/subscribe", strings.NewReader(form.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		w := httptest.NewRecorder()

		h.CreateSubscription(w, req)
		res := w.Result()
		defer res.Body.Close()
		data, err := io.ReadAll(res.Body)

		require.NoError(t, err)
		require.Empty(t, data)

		require.Equal(t, w.Result().StatusCode, http.StatusConflict)
	})
}