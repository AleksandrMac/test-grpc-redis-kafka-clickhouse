package grpc

import (
	"context"
	"log"

	gen "github.com/AleksandrMac/test_hezzl/pkg/user/grpc/userservice"
	"github.com/AleksandrMac/test_hezzl/pkg/user/model"
	"github.com/AleksandrMac/test_hezzl/pkg/user/service"
)

type GRPC struct {
	gen.UnimplementedUserServiceServer
	Services service.Services
}

func (x *GRPC) CreateUser(ctx context.Context, user *gen.User) (*gen.Reply, error) {
	if id, err := x.Services.Create(ctx, user.Email); err != nil {
		return &gen.Reply{Id: "", Status: err.Error()}, err
	} else {
		log.Default().Println("добавлен новый пользователь")
		err = x.Services.LogNewUser(ctx)
		if err != nil {
			// по хорошему в сервисы нужно добавить функции ведения логов,
			// и логировать через них
			// т.к. добавление частичная потеря статистических данных не важна
			// просто выводим сообщение в лог
			log.Default().Println(err)
		}
		return &gen.Reply{Id: id, Status: "OK"}, nil
	}
}
func (x *GRPC) DropUser(ctx context.Context, user *gen.User) (*gen.Reply, error) {
	if err := x.Services.Delete(ctx, user.Id); err != nil {
		return &gen.Reply{Status: err.Error()}, err
	} else {
		return &gen.Reply{Status: "OK"}, nil
	}
}
func (x *GRPC) GetUsers(ctx context.Context, params *gen.SelectParams) (*gen.UserList, error) {
	var usersService []model.User

	usersService, err := x.Services.GetListCache(ctx, params.Limit, params.Offset)

	switch {
	case err != nil:
		return nil, err
	case usersService == nil:
		usersService, err = x.Services.GetList(ctx, params.Limit, params.Offset)
		if err != nil {
			return nil, err
		}
		err = x.Services.SetListCache(ctx, params.Limit, params.Offset, usersService)
		if err != nil {
			log.Default().Println(err)
		}
	}

	users := []*gen.User{}
	for _, userService := range usersService {
		users = append(users, &gen.User{
			Id:    userService.ID,
			Email: userService.Email,
		})
	}
	return &gen.UserList{User: users}, nil
}
