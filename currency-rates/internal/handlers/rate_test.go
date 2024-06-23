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
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type RateHandlerTestSuite struct {
	suite.Suite
	logger *zap.SugaredLogger
}

type successRateClient struct{}

func (c successRateClient) GetDollarRate(ctx context.Context) (float64, error) {
	return 10.0, nil
}

type failedRateClient struct{}

func (c failedRateClient) GetDollarRate(ctx context.Context) (float64, error) {
	return 0.0, errors.New("banks not available")
}

func (t *RateHandlerTestSuite) SetupSuite() {
	l := zap.NewNop()
	t.logger = l.Sugar()
}

func (t *RateHandlerTestSuite) TestSuccessfulGetCurrentRate() {
	expectedRate := float64(10.0)
	s := successRateClient{}
	h := handlers.NewRateHandler(t.logger, s)

	req := httptest.NewRequest(http.MethodGet, "/rates", nil)
	w := httptest.NewRecorder()

	h.GetCurrentRate(w, req)
	res := w.Result()

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	require.NoError(t.T(), err)

	rate, err := strconv.ParseFloat(string(data), 64)
	require.NoError(t.T(), err)
	require.Equal(t.T(), expectedRate, rate)

	require.Equal(t.T(), w.Result().StatusCode, http.StatusOK)
}
func (t *RateHandlerTestSuite) TestFailureGetCurrentRate() {
	f := failedRateClient{}
	h := handlers.NewRateHandler(t.logger, f)

	req := httptest.NewRequest(http.MethodGet, "/rates", nil)
	w := httptest.NewRecorder()

	h.GetCurrentRate(w, req)
	res := w.Result()

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)

	require.NoError(t.T(), err)
	require.Empty(t.T(), data)

	require.Equal(t.T(), w.Result().StatusCode, http.StatusBadRequest)
}

func TestRateHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(RateHandlerTestSuite))
}