package storage

import (
	"context"

	"main.go/clients/events"
)

type Storage interface {
	Save(ctx context.Context, UserData *Userdata) error
	Remove(ctx context.Context, userData *Userdata) error
	Exists(ctx context.Context, userData *Userdata) (bool, error)
	RetrieveLocation(ctx context.Context, userData *Userdata) (*events.CityData, error)
	IsKeyboardMenu(ctx context.Context, userData *Userdata) (bool, error)
}

type Userdata struct {
	UserName string
	ChatId   int
	City     string
	State    string
	Country  string
	Menu     string
}
