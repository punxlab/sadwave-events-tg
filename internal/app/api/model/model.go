package model

import "context"

type SadwaveAPI interface {
	Events(ctx context.Context, city string) ([]*Event, error)
	Cities(ctx context.Context) ([]*City, error)
}

type Date struct {
	Date string `json:"date"`
	Time string `json:"time"`
}

type Event struct {
	Name     string `json:"name"`
	Overview string `json:"overview"`
	Url      string `json:"url"`
	Date     Date   `json:"date"`
}

type City struct {
	Alias string `json:"alias"`
	Name  string `json:"name"`
}
