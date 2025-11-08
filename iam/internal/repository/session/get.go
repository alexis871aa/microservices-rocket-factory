package session

import (
	"context"
	"errors"

	redigo "github.com/gomodule/redigo/redis"

	"github.com/alexis871aa/microservices-rocket-factory/iam/internal/model"
	repoConverter "github.com/alexis871aa/microservices-rocket-factory/iam/internal/repository/converter"
	repoModel "github.com/alexis871aa/microservices-rocket-factory/iam/internal/repository/model"
)

func (r *repository) Get(ctx context.Context, sessionUuid string) (model.Session, model.User, error) {
	cacheKey := r.GetCacheKey(sessionUuid)

	values, err := r.cache.HGetAll(ctx, cacheKey)
	if err != nil {
		if errors.Is(err, redigo.ErrNil) {
			return model.Session{}, model.User{}, model.ErrSessionNotFound
		}
		return model.Session{}, model.User{}, err
	}

	if len(values) == 0 {
		return model.Session{}, model.User{}, model.ErrSessionNotFound
	}

	var sessionRedisView repoModel.SessionRedisView
	err = redigo.ScanStruct(values, &sessionRedisView)
	if err != nil {
		return model.Session{}, model.User{}, err
	}

	session, user := repoConverter.SessionAndUserFromRedisView(sessionRedisView)

	return session, user, nil
}
