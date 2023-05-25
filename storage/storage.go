package storage

import "context"

type Storage interface {
	Save(ctx context.Context, UserData *Userdata) error
	Remove(ctx context.Context, userData Userdata) error
	Exists(ctx context.Context, userData Userdata) (bool, error)
	RetrieveCity(ctx context.Context, userData Userdata) (string, error)
}

type Userdata struct {
	UserName string `json:"username"`
	ChatId   int
	City     string `json:"city"`
}
