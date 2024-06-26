package handlers_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/RinaDish/currency-rates/internal/handlers"
	"github.com/RinaDish/currency-rates/internal/handlers/mocks"
	"github.com/RinaDish/currency-rates/tools"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type SubscribeHandlerTestSuite struct {
	suite.Suite
	logger  tools.Logger
}

func (t *SubscribeHandlerTestSuite) SetupSuite() {
	l := zap.NewNop()
	logger := l.Sugar()

	t.logger = tools.NewZapLogger(logger)
}

func (t *SubscribeHandlerTestSuite)TestSuccessfulCreateSubscription() {	
	mockDB := mocks.NewDb(t.T())
	mockDB.On("SetEmail", mock.Anything, "test@test.com").Return(nil)

	h := handlers.NewSubscribeHandler(t.logger, mockDB)

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
	mockDB := mocks.NewDb(t.T())

	h := handlers.NewSubscribeHandler(t.logger, mockDB)

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
	mockDB := mocks.NewDb(t.T())
	mockDB.On("SetEmail", mock.Anything, "test@gmail.com").Return(errors.New("email exist"))

	h := handlers.NewSubscribeHandler(t.logger, mockDB)

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