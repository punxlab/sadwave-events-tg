package api

import (
	"context"
	"fmt"

	"github.com/punxlab/sadwave-events-tg/internal/app/api/model"
	"github.com/punxlab/sadwave-events-tg/internal/http"
)

type api struct {
	client http.Client
}

func NewSadwaveAPI(client http.Client) model.SadwaveAPI {
	return &api{
		client: client,
	}
}

func (a *api) Events(ctx context.Context, city string) ([]*model.Event, error) {
	res := make([]*model.Event, 0)
	err := a.client.Get(ctx, fmt.Sprintf("/api/events/%s", city), &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *api) Cities(ctx context.Context) ([]*model.City, error) {
	res := make([]*model.City, 0)
	err := a.client.Get(ctx, "/api/cities", &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
