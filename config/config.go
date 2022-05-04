package config

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServiceName string `env:"SERVICE_NAME" envDefault:"TestHEZZL"`

	DBURL string `env:"DB_URL" envDefault:"postgres://testhezzl_user:testpass@192.168.1.13:5432/testhezzl_db?sslmode=disable"`
	// DBConnectTTL - Время для перессоединения, в случае если бд не успела подняться
	DBConnectTTL    uint64 `env:"DB_CONNECT_TTL" envDefault:"5"` // second
	DBMigrateSourse string `env:"DB_MIGRATE_SOURCE" envDefault:"file://db/migration/postgresql"`
	DBMigrateStep   int64  `env:"DB_MIGRATE_STEP" envDefault:"1"`

	ServerGRPCPort uint64 `env:"SERVER_GRPC_PORT" envDefault:"8999"`

	RedisHost     string `env:"REDIS_HOST" envDefault:"localhost"`
	RedisPort     uint64 `env:"REDIS_PORT" envDefault:"6379"`
	RedisCacheTTL uint64 `env:"REDIS_CACHE_TTL" envDefault:"60"` // second

	KafkaBrokers []string `env:"KAFKA_BROKERS" envDefault:"localhost:9092"`
}

func New() (*Config, error) {
	c := &Config{}
	err := env.Parse(c)
	return c, err
}
