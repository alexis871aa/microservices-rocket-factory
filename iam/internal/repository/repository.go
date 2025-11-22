package repository

import (
	"context"
	"time"

	"github.com/alexis871aa/microservices-rocket-factory/iam/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, info model.UserInfo, password string) (string, error)
	GetUser(ctx context.Context, userUUID string) (*model.User, error)
	GetUserByLogin(ctx context.Context, login, password string) (*model.User, error)
}

type SessionRepository interface {
	Create(ctx context.Context, session model.Session, user model.User, ttl time.Duration) error
	Get(ctx context.Context, sessionUuid string) (model.Session, model.User, error)
	AddSessionToUserSet(ctx context.Context, userUuid, sessionUuid string) error
}
