package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RinaDish/subscription-sender/internal/handlers"
	"github.com/RinaDish/subscription-sender/internal/handlers/mocks"
	"github.com/RinaDish/subscription-sender/internal/services"
	"github.com/RinaDish/subscription-sender/tools"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type NotifyHandlerTestSuite struct {
	suite.Suite
	logger  tools.Logger
}

func (t *NotifyHandlerTestSuite) SetupSuite() {
	l := zap.NewNop()
	logger := l.Sugar()

	t.logger = tools.NewZapLogger(logger)
}

func (t *NotifyHandlerTestSuite) TestSuccessfulNotifySubscribers() {
	notification := handlers.Notification{
		Rate:   10.09,
		Emails: []string{"test@test.com", "testtest@test.com"},
	}

	body, _ := json.Marshal(notification)

	mockNotifyService := mocks.NewNotifyService(t.T())
	handler := handlers.NewNotifyHandler(t.logger, mockNotifyService)

	req := httptest.NewRequest(http.MethodPost, "/notify", bytes.NewReader(body))
	w := httptest.NewRecorder()

	mockNotifyService.On("NotifySubscribers", mock.Anything, services.Notification{Rate: notification.Rate, Emails: notification.Emails}).Once()

	handler.NotifySubscribers(w, req)

	res := w.Result()
	defer res.Body.Close()

	require.Equal(t.T(), http.StatusOK, res.StatusCode)
}

func (t *NotifyHandlerTestSuite) TestFailureInvalidJsonBody() {
	body := []byte("invalid body")

	mockNotifyService := mocks.NewNotifyService(t.T())
	handler := handlers.NewNotifyHandler(t.logger, mockNotifyService)

	req := httptest.NewRequest(http.MethodPost, "/notify", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.NotifySubscribers(w, req)

	res := w.Result()
	defer res.Body.Close()

	require.Equal(t.T(), http.StatusBadRequest, res.StatusCode)
}

func TestNotifyHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(NotifyHandlerTestSuite))
}