package session

import (
	"fmt"

	def "github.com/alexis871aa/microservices-rocket-factory/iam/internal/service"
	"github.com/alexis871aa/microservices-rocket-factory/platform/pkg/cache"
)

var _ def.SessionRepository = (*repository)(nil)

const (
	cacheKeyPrefix = "iam:session:"
)

type repository struct {
	cache cache.RedisClient
}

func NewRepository(cache cache.RedisClient) *repository {
	return &repository{
		cache: cache,
	}
}

func (r *repository) GetCacheKey(uuid string) string {
	return fmt.Sprintf("%s:%s", cacheKeyPrefix, uuid)
}
