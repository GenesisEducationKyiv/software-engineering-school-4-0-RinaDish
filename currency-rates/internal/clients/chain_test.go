package clients_test

import (
	"context"
	"errors"
	"testing"

	"github.com/RinaDish/currency-rates/internal/clients"
	"github.com/RinaDish/currency-rates/internal/clients/mocks"
	"github.com/RinaDish/currency-rates/tools"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

var expectedSuccessRate = float64(10.0)
var expectedFailRate = float64(0.0)

type ChainTestSuite struct {
	suite.Suite
	logger        tools.Logger
	successClient *mocks.RateClient
	failClient *mocks.RateClient
}

func (t *ChainTestSuite) SetupSuite() {
	l := zap.NewNop()
	logger := l.Sugar()

	t.logger = tools.NewZapLogger(logger)

	mockRateClientS := mocks.NewRateClient(t.T())
	mockRateClientS.On("GetDollarRate", mock.Anything).Return(expectedSuccessRate, nil)
	t.successClient = mockRateClientS

	mockRateClientF := mocks.NewRateClient(t.T())
	mockRateClientF.On("GetDollarRate", mock.Anything).Return(expectedFailRate, errors.New("bank not available"))
	t.failClient = mockRateClientF
}

func setChain(cf clients.RateClient, cs clients.RateClient) *clients.BaseChain {
	clf := clients.NewBaseChain(cf)
	cls := clients.NewBaseChain(cs)
	clf.SetNext(cls)

	return clf
}

func (t *ChainTestSuite) TestSuccessFirstClient() {
	client := setChain(t.successClient, t.successClient)

	rate, err := client.GetDollarRate(context.Background())

	require.NoError(t.T(), err)
	require.Equal(t.T(), expectedSuccessRate, rate)

	t.successClient.AssertExpectations(t.T())
}

func (t *ChainTestSuite) TestSuccessSecondClient() {
	client := setChain(t.failClient, t.successClient)

	rate, err := client.GetDollarRate(context.Background())

	require.NoError(t.T(), err)
	require.Equal(t.T(), expectedSuccessRate, rate)

	t.failClient.AssertExpectations(t.T())
	t.successClient.AssertExpectations(t.T())
}

func (t *ChainTestSuite) TestFailClient() {
	client := setChain(t.failClient, t.failClient)

	rate, err := client.GetDollarRate(context.Background())

	require.Error(t.T(), err)
	require.Equal(t.T(), expectedFailRate, rate)
	require.Contains(t.T(), err.Error(), "bank not available")

	t.failClient.AssertExpectations(t.T())
}

func TestChainTestSuite(t *testing.T) {
	suite.Run(t, new(ChainTestSuite))
}