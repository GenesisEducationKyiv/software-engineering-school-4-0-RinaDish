package clients

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

  "github.com/RinaDish/currency-rates/tools"
)

type PrivatRate struct {
	Sale	string 	`json:"sale"`
  Ccy string `json:"ccy"`
}

type PrivatClient struct {
	logger tools.Logger
	client *http.Client
}

func NewPrivatClient(logger tools.Logger) *PrivatClient {
  return &PrivatClient{
		logger: logger.With("client", "PrivatBank"),
		client: http.DefaultClient,
	}
}

func (privatClient *PrivatClient) GetDollarRate(ctx context.Context) (float64, error)  {
  ctx, cancel := context.WithTimeout(ctx, ApiCallTimeout)
	defer cancel()

  url := "https://api.privatbank.ua/p24api/pubinfo?json&exchange&coursid=5"

  req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

  if err != nil {
    return 0.0, err
  }
  
  res, err := privatClient.client.Do(req)
  if err != nil {
    return 0.0, err
  }

  defer res.Body.Close()

  body, err := io.ReadAll(res.Body)
  if err != nil {
    return 0.0, err
  }

  var ans []PrivatRate
  err = json.Unmarshal(body, &ans)
  if err != nil {
	  return 0.0, err
  }

  privatClient.logger.Info("Rate: ", ans)

  if len(ans) > 0 {
    for _, val := range ans {
      if val.Ccy == "USD" {
        return strconv.ParseFloat(val.Sale, 64)
      }
    }
  }
  
  return 0.0, errors.New("currency not found")
}