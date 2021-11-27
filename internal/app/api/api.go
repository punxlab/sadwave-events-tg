package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/punxlab/sadwave-events-tg/internal/app/api/model"
	"net/http"
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
	resp, err := a.client.
		R().
		SetContext(ctx).
		Get(fmt.Sprintf("events/%s", city))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() == http.StatusNotFound {
		return nil, ErrNotFound
	}

	if !resp.IsSuccess() {
		return nil, errors.Errorf("get city: unexpected status %s: %s", resp.Status(), string(resp.Body()))
	}

	res := make([]*model.Event, 0)
	err = json.Unmarshal(resp.Body(), &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *api) Cities(ctx context.Context) ([]*model.City, error) {
	res := make([]*model.City, 0)
	resp, err := a.client.
		R().
		SetContext(ctx).
		Get("/cities")
	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, errors.Errorf("get city: unexpected status %s: %s", resp.Status(), string(resp.Body()))
	}

	err = json.Unmarshal(resp.Body(), &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
