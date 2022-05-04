package user

import (
	"time"

	"github.com/AleksandrMac/test_hezzl/pkg/user/cache"
	"github.com/AleksandrMac/test_hezzl/pkg/user/logger"
	"github.com/AleksandrMac/test_hezzl/pkg/user/service"
	"github.com/AleksandrMac/test_hezzl/pkg/user/storage"
	"github.com/AleksandrMac/test_hezzl/pkg/user/storage/pgxdb"

	userRedis "github.com/AleksandrMac/test_hezzl/pkg/user/cache/redis"
	userGRPC "github.com/AleksandrMac/test_hezzl/pkg/user/grpc"
	userGRPCgen "github.com/AleksandrMac/test_hezzl/pkg/user/grpc/userservice"
	"github.com/AleksandrMac/test_hezzl/pkg/user/logger/kafka"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"

	"google.golang.org/grpc"
)

// NewService необходимый функционал для работы с сущностью user.
func NewService(storage storage.CRUDL, log logger.Logger, cache cache.Cacher) service.Services {
	return service.New(storage, log, cache)
}

func NewStoragePGX(dbpool *pgxpool.Pool) storage.CRUDL {
	return pgxdb.New(dbpool)
}

func NewLoggerKafka(kafkaBrokers []string) (*kafka.Kafka, error) {
	return kafka.New(kafkaBrokers)
}

func NewCacherRedis(client *redis.Client, ttl time.Duration) cache.Cacher {
	return userRedis.New(client, ttl)
}

func RegistrationGRPCService(s grpc.ServiceRegistrar, service service.Services) *userGRPC.GRPC {
	g := &userGRPC.GRPC{
		Services: service,
	}
	userGRPCgen.RegisterUserServiceServer(s, g)
	return g
}
