package model

import "context"

type SadwaveAPI interface {
	Events(ctx context.Context, city string) ([]*Event, error)
	Cities(ctx context.Context) ([]*City, error)
}

type Event struct {
	Title           string
	DescriptionHTML string
	ImageURL        string
}

type City struct {
	Code string
	Name string
}
