package handlers

import (
	"context"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

type RateClient interface {
	GetDollarRate(ctx context.Context) (float64, error) 
}

type RateHandler struct {
	l   *zap.SugaredLogger
	rateClient RateClient
}

func NewRateHandler(l *zap.SugaredLogger, r RateClient) RateHandler {

	return RateHandler{
		l: l,
		rateClient: r,
	}
}

func (h RateHandler) GetCurrentRate(w http.ResponseWriter, r *http.Request) {
	rate, err := h.rateClient.GetDollarRate(context.Background())

	w.Header().Set("Content-Type", "application/json")
	if err == nil { 
		w.WriteHeader(http.StatusOK)
		strRate := strconv.FormatFloat(rate, 'f', -1, 64)
		_, err := w.Write([]byte(strRate))
		if err != nil {
			h.l.Error(err)
		}
		return
	} 
	w.WriteHeader(http.StatusBadRequest)
	
}