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
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type SubscribeHandlerTestSuite struct {
	suite.Suite
	logger *zap.SugaredLogger
}

type successDb struct{}

func (d successDb) SetEmail(ctx context.Context, email string) error {
	return nil
}

type failDb struct{}

func (d failDb) SetEmail(ctx context.Context, email string) error {
	return errors.New("email exist")
}

func (t *SubscribeHandlerTestSuite) SetupSuite() {
	l := zap.NewNop()
	t.logger = l.Sugar()
}

func (t *SubscribeHandlerTestSuite)TestSuccessfulCreateSubscription() {	
	d := successDb{}
	h := handlers.NewSubscribeHandler(t.logger, d)

	form := url.Values{}
	form.Add("email", "test@test.com")

	req := httptest.NewRequest(http.MethodPost, "/subscribe", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	h.CreateSubscription(w, req)
	res := w.Result()

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)

	require.NoError(t.T(), err)
	require.Empty(t.T(), data)

	require.Equal(t.T(), w.Result().StatusCode, http.StatusOK)
}

func (t *SubscribeHandlerTestSuite)TestFailureInvalidEmail() {	
	d := successDb{}
	h := handlers.NewSubscribeHandler(t.logger, d)

	form := url.Values{}
	form.Add("email", "test")

	req := httptest.NewRequest(http.MethodPost, "/subscribe", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	h.CreateSubscription(w, req)
	res := w.Result()

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)

	require.NoError(t.T(), err)
	require.Equal(t.T(), "Invalid email\n", string(data))

	require.Equal(t.T(), w.Result().StatusCode, http.StatusConflict)
}

func (t *SubscribeHandlerTestSuite)TestFailureDbSetEmail() {	
	d := failDb{}
	h := handlers.NewSubscribeHandler(t.logger, d)

	form := url.Values{}
	form.Add("email", "test@gmail.com")

	req := httptest.NewRequest(http.MethodPost, "/subscribe", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	h.CreateSubscription(w, req)
	res := w.Result()
	
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)

	require.NoError(t.T(), err)
	require.Empty(t.T(), data)

	require.Equal(t.T(), w.Result().StatusCode, http.StatusConflict)
}

func TestSubscribeHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(SubscribeHandlerTestSuite))
}