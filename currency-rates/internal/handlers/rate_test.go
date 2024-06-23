package handlers_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/RinaDish/currency-rates/internal/handlers"
	"github.com/RinaDish/currency-rates/internal/handlers/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type RateHandlerTestSuite struct {
	suite.Suite
	logger *zap.SugaredLogger
}

func (t *RateHandlerTestSuite) SetupSuite() {
	l := zap.NewNop()
	t.logger = l.Sugar()
}

func (t *RateHandlerTestSuite) TestSuccessfulGetCurrentRate() {
	expectedRate := float64(10.0)

	mockRateClient := mocks.NewRateClient(t.T())
	mockRateClient.On("GetDollarRate", mock.Anything).Return(expectedRate, nil)
	
	h := handlers.NewRateHandler(t.logger, mockRateClient)

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
	mockRateClient := mocks.NewRateClient(t.T())
	mockRateClient.On("GetDollarRate", mock.Anything).Return( 0.0, errors.New("banks not available"))

	h := handlers.NewRateHandler(t.logger, mockRateClient)

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