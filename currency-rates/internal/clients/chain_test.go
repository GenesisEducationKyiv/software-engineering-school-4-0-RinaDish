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
var errBankNotAvailable = errors.New("bank not available")

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
	mockRateClientF.On("GetDollarRate", mock.Anything).Return(expectedFailRate, errBankNotAvailable)
	t.failClient = mockRateClientF
}

func setChain(clientFirst clients.RateClient, clientSecond clients.RateClient) *clients.BaseChain {
	clf := clients.NewBaseChain(clientFirst)
	cls := clients.NewBaseChain(clientSecond)
	clf.SetNext(cls)

	return clf
}

func (t *ChainTestSuite) TestSuccessFirstClient() {
	client := setChain(t.successClient, t.successClient)

	actualRate, err := client.GetDollarRate(context.Background())

	require.NoError(t.T(), err)
	require.Equal(t.T(), expectedSuccessRate, actualRate)
}

func (t *ChainTestSuite) TestSuccessSecondClient() {
	client := setChain(t.failClient, t.successClient)

	actualRate, err := client.GetDollarRate(context.Background())

	require.NoError(t.T(), err)
	require.Equal(t.T(), expectedSuccessRate, actualRate)
}

func (t *ChainTestSuite) TestFailClient() {
	client := setChain(t.failClient, t.failClient)

	actualRate, err := client.GetDollarRate(context.Background())

	require.Error(t.T(), err)
	require.Equal(t.T(), expectedFailRate, actualRate)
	require.True(t.T(), errors.Is(err, errBankNotAvailable))
}

func TestChainTestSuite(t *testing.T) {
	suite.Run(t, new(ChainTestSuite))
}
