package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/RinaDish/currency-rates/tools"
)

type RateClient interface {
	GetDollarRate(ctx context.Context) (float64, error) 
}

type RateHandler struct {
	logger   tools.Logger
	rateClient RateClient
}

func NewRateHandler(logger tools.Logger, client RateClient) RateHandler {
	return RateHandler{
		logger: logger,
		rateClient: client,
	}
}

func (hadler RateHandler) GetCurrentRate(w http.ResponseWriter, r *http.Request) {
	rate, err := hadler.rateClient.GetDollarRate(context.Background())

	w.Header().Set("Content-Type", "application/json")
	if err == nil { 
		w.WriteHeader(http.StatusOK)

		strRate := strconv.FormatFloat(rate, 'f', -1, 64)
		_, err := w.Write([]byte(strRate))
		if err != nil {
			hadler.logger.Error(err)
		}

		return
	} 
	
	w.WriteHeader(http.StatusBadRequest)
}