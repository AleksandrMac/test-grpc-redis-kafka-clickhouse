package cache

import (
	"context"

	"github.com/AleksandrMac/test_hezzl/pkg/user/model"
)

type Cacher interface {
	GetList(ctx context.Context, limit, offset uint64) ([]model.User, error)
	SetList(ctx context.Context, limit, offset uint64, users []model.User) error
}

type Cache struct {
	Cacher
}
