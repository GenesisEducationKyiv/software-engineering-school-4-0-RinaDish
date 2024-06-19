package clients

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"
)

type NBURate struct {
	Rate	float64 	`json:"rate"`
}

type NBUClient struct {
	logger *zap.SugaredLogger
	client *http.Client
}

func NewNBUClient(logger *zap.SugaredLogger) NBUClient {
	client := &http.Client {
	}

	return NBUClient{
		logger: logger.With("client", "NBU"),
		client: client,
	}
}

func (nbuClient NBUClient) GetDollarRate(ctx context.Context) (float64, error){
  url := "https://bank.gov.ua/NBUStatService/v1/statdirectory/dollar_info?json"

  req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

  if err != nil {
    nbuClient.logger.Info(err)
    return 0.0, err
  }

  res, err := nbuClient.client.Do(req)

  if err != nil {
    nbuClient.logger.Info(err)
    return 0.0, err
  }

  defer res.Body.Close()

  body, err := io.ReadAll(res.Body)

  if err != nil {
    nbuClient.logger.Info(err)
    return 0.0, err
  }

  var ans []NBURate
  err = json.Unmarshal(body, &ans)

  if err != nil {
    nbuClient.logger.Info(err)
    return 0.0, err
  }

  nbuClient.logger.Info("Rate: ", ans[0].Rate)
  
  return ans[0].Rate, nil
}