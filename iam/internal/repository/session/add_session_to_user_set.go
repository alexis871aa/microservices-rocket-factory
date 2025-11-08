package session

import "context"

func (r *repository) AddSessionToUserSet(ctx context.Context, userUuid, sessionUuid string) error {
	cacheKey := r.GetCacheKey(userUuid)
	return r.cache.SAdd(ctx, cacheKey, sessionUuid)
}
