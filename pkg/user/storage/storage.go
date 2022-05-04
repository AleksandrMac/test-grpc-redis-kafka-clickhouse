package storage

import (
	"context"

	"github.com/AleksandrMac/test_hezzl/pkg/user/model"
)

type Storage struct {
	DB CRUDL
}

type CRUDL interface {
	Create(ctx context.Context, email string) (id uint64, err error)
	Delete(ctx context.Context, id uint64) error
	GetList(ctx context.Context, limit, offset uint64) ([]model.UserDB, error)
}
