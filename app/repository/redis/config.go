package redisrepo

import (
	"time"
)

func (r *redisRepo) Enabled() bool {
	return r.UseRedis
}

func (r *redisRepo) GetTTL() time.Duration {
	return r.DefaultTTL
}
