package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/punxlab/sadwave-events-tg/internal/app/api/model"
)

var ErrNotFound = errors.New("not found")

type api struct {
	client *resty.Client
}

func NewSadwaveAPI(apiUrl string) model.SadwaveAPI {
	c := resty.New()
	c.SetHostURL(apiUrl)

	return &api{
		client: c,
	}
}

func (a *api) Events(ctx context.Context, city string) ([]*model.Event, error) {
	res := make([]*model.Event, 0)
	err := a.do(func() (*resty.Response, error) {
		return a.client.
			R().
			SetContext(ctx).
			SetResult(&res).
			Get(fmt.Sprintf("events/%s", city))
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *api) Cities(ctx context.Context) ([]*model.City, error) {
	res := make([]*model.City, 0)
	err := a.do(func() (*resty.Response, error) {
		return a.client.
			R().
			SetContext(ctx).
			SetResult(&res).
			Get("/cities")
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *api) do(f func() (*resty.Response, error)) error {
	resp, err := f()
	if err != nil {
		return err
	}

	if resp.StatusCode() == http.StatusNotFound {
		return ErrNotFound
	}

	if !resp.IsSuccess() {
		return errors.Errorf("unexpected status %s: %s", resp.Status(), string(resp.Body()))
	}

	return nil
}
