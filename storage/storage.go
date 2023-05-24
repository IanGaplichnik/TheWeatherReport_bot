package storage

import "context"

type Storage interface {
	Save(ctx context.Context, UserData *Userdata) error
	Remove(ctx context.Context, userData Userdata) error
	Exists(ctx context.Context, userData Userdata) (bool, error)
}

type Userdata struct {
	UserName string `json:"username"`
	City     string `json:"city"`
}
