package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/go-redis/redis/v8"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/AleksandrMac/test_hezzl/config"
	"github.com/AleksandrMac/test_hezzl/pkg/user"
)

type DBOptions struct {
	// URL адрес подключения к бд
	URL string
	// ConnectTTL время допустимое для подключения к бд
	ConnectTTL uint64
	// MigrateSource путь с файлами для миграций
	MigrateSource string
	// MigrateStep количество шагов миграций
	MigrateStep int64
}

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.ServerGRPCPort))
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}

	// register server
	grpcServer := grpc.NewServer()
	ctxMain, cancelMain := context.WithCancel(context.Background())
	defer cancelMain()

	pgDB, err := NewDBPool(ctxMain, &DBOptions{
		URL:           cfg.DBURL,
		ConnectTTL:    cfg.DBConnectTTL,
		MigrateSource: cfg.DBMigrateSourse,
		MigrateStep:   cfg.DBMigrateStep,
	})
	if err != nil {
		log.Panic(err)
	}
	log.Default().Println("PostgreSQL connected")
	defer pgDB.Close()

	redisClient, err := NewRedisClient(cfg.RedisHost, cfg.RedisPort)
	if err != nil {
		log.Panic(err)
	}
	log.Default().Println("Redis connected")
	defer redisClient.Close()

	kafkaBrokers, err := user.NewLoggerKafka(cfg.KafkaBrokers)
	if err != nil {
		log.Panic(err)
	}
	log.Default().Println("Kafka connected")
	defer kafkaBrokers.Shutdown(ctxMain)

	userService := user.NewService(
		user.NewStoragePGX(pgDB),
		kafkaBrokers,
		user.NewCacherRedis(redisClient, time.Duration(cfg.RedisCacheTTL)*time.Second),
	)

	user.RegistrationGRPCService(grpcServer, userService)

	reflection.Register(grpcServer)

	// start to serve
	log.Default().Println("GRPC server listen")
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

func NewDBPool(ctx context.Context, opts *DBOptions) (dbpool *pgxpool.Pool, err error) {

	url, ttl, migrationPath := opts.URL, opts.ConnectTTL, opts.MigrateSource
	for ; ttl > 0; ttl-- {
		dbpool, err = pgxpool.Connect(ctx, url)
		if err == nil {
			break
		}
		log.Default().Printf("Неудачная попытка подключения к БД. Осталось попыток %d", ttl-1)
		<-time.After(time.Second)
	}
	if err != nil {
		return nil, err
	}

	if err = dbpool.Ping(ctx); err != nil {
		return nil, err
	}

	m, err := migrate.New(migrationPath, url)
	if err != nil {
		return nil, fmt.Errorf("NewDBPool:CreateMigration:%w", err)
	}
	if err := m.Steps(int(opts.MigrateStep)); err != nil {
		if err == fs.ErrNotExist {
			log.Default().Println("Новых миграций не обнаружено")
			return dbpool, nil
		}
		return nil, fmt.Errorf("NewDBPool:MakeMigration:%w", err)
	}
	return dbpool, nil
}

func NewRedisClient(host string, port uint64) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, fmt.Errorf("try to ping to redis: %w", err)
	}

	return client, nil
}
