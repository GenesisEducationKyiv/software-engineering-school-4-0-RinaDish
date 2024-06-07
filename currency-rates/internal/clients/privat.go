package clients

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

type PrivatRate struct {
	Sale	string 	`json:"sale"`
  Ccy string `json:"ccy"`
}

type PrivatClient struct {
	l *zap.SugaredLogger
	c *http.Client
}

func NewPrivatClient(l *zap.SugaredLogger) PrivatClient {
	client := &http.Client {
	}

	return PrivatClient{
		l: l.With("client", "PrivatBank"),
		c: client,
	}
}

func (n PrivatClient) GetDollarRate(ctx context.Context) (float64, error)  {


  url := "https://api.privatbank.ua/p24api/pubinfo?json&exchange&coursid=5"

  req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

  if err != nil {
    n.l.Info(err)
    return 0.0, err
  }
  
  res, err := n.c.Do(req)
  if err != nil {
    n.l.Info(err)
    return 0.0, err
  }
  defer res.Body.Close()

  body, err := io.ReadAll(res.Body)
  if err != nil {
    n.l.Info(err)
    return 0.0, err
  }

  var ans []PrivatRate
  err = json.Unmarshal(body, &ans)
  if err != nil {
	n.l.Info(err)
	return 0.0, err
  }

  for _, val := range ans {
    if val.Ccy == "USD" {
      return strconv.ParseFloat(val.Sale, 64)
    }
  }

  return 0.0, errors.New("currency not found")
}