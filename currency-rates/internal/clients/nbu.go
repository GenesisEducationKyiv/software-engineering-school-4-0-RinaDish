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
	l *zap.SugaredLogger
	c *http.Client
}

func NewNBUClient(l *zap.SugaredLogger) NBUClient {
	client := &http.Client {
	}

	return NBUClient{
		l: l.With("client", "NBU"),
		c: client,
	}
}

func (n NBUClient) GetDollarRate(ctx context.Context) (float64, error)  {


  url := "https://bank.gov.ua/NBUStatService/v1/statdirectory/dollar_info?json"

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

  var ans []NBURate
  err = json.Unmarshal(body, &ans)
  if err != nil {
	n.l.Info(err)
	return 0.0, err
  }

  n.l.Info("Rate: ", ans[0].Rate)
  return ans[0].Rate, nil
}