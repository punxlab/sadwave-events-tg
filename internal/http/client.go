package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Client interface {
	Get(ctx context.Context, action string, res interface{}) error
}

type client struct {
	host   string
	client *http.Client
}

func NewClient(host string, timeout time.Duration) Client {
	return &client{
		host: host,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *client) Get(ctx context.Context, action string, res interface{}) error {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", c.host, action), nil)
	if err != nil {
		return fmt.Errorf("create request: %v", err)
	}

	req = req.WithContext(ctx)
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %v", err)
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response code %v", resp.StatusCode)
	}

	if err := json.Unmarshal(b, res); err != nil {
		return fmt.Errorf("unmarshal response: %v", resp)
	}

	return nil
}
