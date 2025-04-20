package redisrepo

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

func (r *redisRepo) Enabled() bool {
	return r.UseRedis
}

func (r *redisRepo) GetTTL() time.Duration {
	return r.DefaultTTL
}

func (r *redisRepo) Get(ctx context.Context, key string) (value []byte, err error) {
	if value, err = r.Conn.Get(ctx, r.Prefix+key).Bytes(); err != nil {
		logrus.Error("Redis Get:", err)
		return
	}

	return
}

func (r *redisRepo) Set(ctx context.Context, key string, value []byte, expiration *time.Duration) (err error) {
	if expiration == nil {
		defaultTTL := r.GetTTL()
		expiration = &defaultTTL
	}
	if res := r.Conn.Set(ctx, r.Prefix+key, value, *expiration); res.Err() != nil {
		logrus.Error("Redis Set:", err)
		return
	}

	return
}
