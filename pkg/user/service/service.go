package service

import (
	"context"
	"strconv"

	"github.com/AleksandrMac/test_hezzl/pkg/user/cache"
	"github.com/AleksandrMac/test_hezzl/pkg/user/logger"
	"github.com/AleksandrMac/test_hezzl/pkg/user/model"
	"github.com/AleksandrMac/test_hezzl/pkg/user/storage"
)

type Services interface {
	// Create создание пользователя в бд
	Create(ctx context.Context, email string) (id string, err error)
	// Delete явное удаление пользователя из бд
	Delete(ctx context.Context, userID string) error
	// GetList получает список пользователей из БД
	GetList(ctx context.Context, limit, offset uint64) ([]model.User, error)
	// GetListCache получает список пользователей из кэша
	GetListCache(ctx context.Context, limit, offset uint64) ([]model.User, error)
	// SetListCache добавляет список пользователей в кэш
	SetListCache(ctx context.Context, limit, offset uint64, users []model.User) error
	// LogNewUser отправляет уведомление о создании нового пользователя.(в кафку)
	LogNewUser(context.Context) error
}

type Service struct {
	storage storage.CRUDL
	log     logger.Logger
	cache   cache.Cacher
}

func New(storage storage.CRUDL, log logger.Logger, c cache.Cacher) *Service {
	return &Service{
		storage: storage,
		log:     log,
		cache:   c,
	}
}

// Create создание пользователя в бд
func (x *Service) Create(ctx context.Context, email string) (id string, err error) {
	iduint, err := x.storage.Create(ctx, email)
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(iduint, 10), nil
}

// Delete явное удаление пользователя из бд
func (x *Service) Delete(ctx context.Context, userID string) error {
	iduint, err := strconv.ParseUint(userID, 10, 0)
	if err != nil {
		return err
	}
	return x.storage.Delete(ctx, iduint)
}

// GetList получает список пользователей из БД
func (x *Service) GetList(ctx context.Context, limit, offset uint64) ([]model.User, error) {
	usersDB, err := x.storage.GetList(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	users := []model.User{}
	for _, userDB := range usersDB {
		users = append(users, model.User{
			ID:    strconv.FormatUint(userDB.ID, 10),
			Email: userDB.Email,
		})
	}
	return users, nil
}

// GetListCache получает список пользователей из кэша
func (x *Service) GetListCache(ctx context.Context, limit, offset uint64) ([]model.User, error) {
	return x.cache.GetList(ctx, limit, offset)
}

// SetListCache добавляет список пользователей в кэш
func (x *Service) SetListCache(ctx context.Context, limit, offset uint64, users []model.User) error {
	return x.cache.SetList(ctx, limit, offset, users)
}

// LogNewUser отправляет уведомление о создании нового пользователя.(в кафку)
func (x *Service) LogNewUser(ctx context.Context) error {
	return x.log.LogNewUser(ctx)
}
