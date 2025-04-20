package usecase_member

import (
	"app/domain/model"
	"context"
	"encoding/json"
	"time"

	"github.com/sirupsen/logrus"
)

func (u *appUsecase) _CacheConfig(ctx context.Context) (config model.Config) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if u.redisRepo.Enabled() {
		if strval, err := u.redisRepo.Get(ctx, "config"); err == nil {
			if err = json.Unmarshal(strval, &config); err == nil {
				// from cache
				return
			}
		}
	}

	cfg, err := u.mongodbRepo.FetchOneConfig(ctx, map[string]interface{}{})
	if err != nil || cfg == nil {
		logrus.Error("no config available")
	}

	config = *cfg

	// set cache
	if u.redisRepo.Enabled() {
		ttl := 168 * time.Hour

		byteData, _ := json.Marshal(config)
		u.redisRepo.Set(ctx, "config", byteData, &ttl)
	}

	return
}
