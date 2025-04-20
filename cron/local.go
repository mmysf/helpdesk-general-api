package cronjob

import (
	"app/domain/model"
	"context"
	"encoding/json"
	"time"

	"github.com/sirupsen/logrus"
)

func (cj *cronjob) _CacheConfig(ctx context.Context) (config model.Config) {
	_, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	if cj.redisRepo.Enabled() {
		if strval, err := cj.redisRepo.Get(cj.ctx, "config"); err == nil {
			if err = json.Unmarshal(strval, &config); err == nil {
				// from cache
				return
			}
		}
	}

	cfg, err := cj.mongodbRepo.FetchOneConfig(cj.ctx, map[string]interface{}{})
	if err != nil || cfg == nil {
		logrus.Error("no config available")
	}

	config = *cfg

	// set cache
	if cj.redisRepo.Enabled() {
		ttl := 168 * time.Hour

		byteData, _ := json.Marshal(config)
		cj.redisRepo.Set(cj.ctx, "config", byteData, &ttl)
	}

	return
}
