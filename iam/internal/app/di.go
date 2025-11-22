package app

import (
	"context"
	"database/sql"
	"fmt"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"

	authV1API "github.com/alexis871aa/microservices-rocket-factory/iam/internal/api/auth/v1"
	userV1API "github.com/alexis871aa/microservices-rocket-factory/iam/internal/api/user/v1"
	"github.com/alexis871aa/microservices-rocket-factory/iam/internal/config"
	"github.com/alexis871aa/microservices-rocket-factory/iam/internal/migrator"
	"github.com/alexis871aa/microservices-rocket-factory/iam/internal/repository"
	sessionRepository "github.com/alexis871aa/microservices-rocket-factory/iam/internal/repository/session"
	userRepository "github.com/alexis871aa/microservices-rocket-factory/iam/internal/repository/user"
	"github.com/alexis871aa/microservices-rocket-factory/iam/internal/service"
	authService "github.com/alexis871aa/microservices-rocket-factory/iam/internal/service/auth"
	userService "github.com/alexis871aa/microservices-rocket-factory/iam/internal/service/user"
	"github.com/alexis871aa/microservices-rocket-factory/platform/pkg/cache"
	"github.com/alexis871aa/microservices-rocket-factory/platform/pkg/cache/redis"
	"github.com/alexis871aa/microservices-rocket-factory/platform/pkg/closer"
	"github.com/alexis871aa/microservices-rocket-factory/platform/pkg/logger"
	authV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/proto/auth/v1"
	userV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/proto/user/v1"
)

type diContainer struct {
	userV1API userV1.UserServiceServer
	authV1API authV1.AuthServiceServer

	userService service.UserService
	authService service.AuthService

	userRepository    repository.UserRepository
	sessionRepository repository.SessionRepository

	sqlDb   *sql.DB
	pgxConn *pgx.Conn

	migratorRunner *migrator.Migrator

	redisPool   *redigo.Pool
	redisClient cache.RedisClient
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) UserV1API(ctx context.Context) userV1.UserServiceServer {
	if d.userV1API == nil {
		d.userV1API = userV1API.NewAPI(d.UserService(ctx))
	}
	return d.userV1API
}

func (d *diContainer) UserService(ctx context.Context) service.UserService {
	if d.userService == nil {
		d.userService = userService.NewService(d.UserRepository(ctx))
	}
	return d.userService
}

func (d *diContainer) UserRepository(ctx context.Context) repository.UserRepository {
	if d.userRepository == nil {
		d.userRepository = userRepository.NewRepository(d.SqlDB(ctx))
	}
	return d.userRepository
}

func (d *diContainer) AuthV1API(ctx context.Context) authV1.AuthServiceServer {
	if d.authV1API == nil {
		d.authV1API = authV1API.NewAPI(d.AuthService(ctx))
	}

	return d.authV1API
}

func (d *diContainer) AuthService(ctx context.Context) service.AuthService {
	if d.authService == nil {
		d.authService = authService.NewService(d.SessionRepository(ctx), d.UserRepository(ctx), config.AppConfig().Redis.CacheTTL())
	}

	return d.authService
}

func (d *diContainer) SessionRepository(_ context.Context) repository.SessionRepository {
	if d.sessionRepository == nil {
		d.sessionRepository = sessionRepository.NewRepository(d.RedisClient())
	}

	return d.sessionRepository
}

func (d *diContainer) RedisPool() *redigo.Pool {
	if d.redisPool == nil {
		d.redisPool = &redigo.Pool{
			MaxIdle:     config.AppConfig().Redis.MaxIdle(),
			IdleTimeout: config.AppConfig().Redis.IdleTimeout(),
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(ctx, "tcp", config.AppConfig().Redis.Address())
			},
		}

		closer.AddNamed("Redis pool", func(ctx context.Context) error {
			return d.redisPool.Close()
		})
	}

	return d.redisPool
}

func (d *diContainer) RedisClient() cache.RedisClient {
	if d.redisClient == nil {
		d.redisClient = redis.NewClient(d.RedisPool(), logger.Logger(), config.AppConfig().Redis.ConnectionTimeout())
	}

	return d.redisClient
}

func (d *diContainer) SqlDB(ctx context.Context) *sql.DB {
	if d.sqlDb == nil {
		d.sqlDb = stdlib.OpenDB(*d.PgxConn(ctx).Config().Copy())
	}

	return d.sqlDb
}

func (d *diContainer) PgxConn(ctx context.Context) *pgx.Conn {
	if d.pgxConn == nil {
		conn, err := pgx.Connect(ctx, config.AppConfig().Postgres.URI())
		if err != nil {
			panic(fmt.Sprintf("💥 failed to connect to database: %v", err))
		}

		err = conn.Ping(ctx)
		if err != nil {
			panic(fmt.Sprintf("💥 failed to ping database: %v", err))
		}

		closer.AddNamed("PostgresSQL connection", func(ctx context.Context) error {
			return conn.Close(ctx)
		})

		d.pgxConn = conn
	}

	return d.pgxConn
}

func (d *diContainer) MigratorRunner(ctx context.Context) *migrator.Migrator {
	if d.migratorRunner == nil {
		d.migratorRunner = migrator.NewMigrator(
			d.SqlDB(ctx),
			config.AppConfig().Postgres.MigrationsDir(),
		)
	}
	return d.migratorRunner
}
