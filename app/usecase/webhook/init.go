package usecase_webhook

import (
	mongorepo "app/app/repository/mongo"
	redisrepo "app/app/repository/redis"
	"context"
	"time"

	"github.com/Yureka-Teknologi-Cipta/yureka/response"
)

type webhookUsecase struct {
	mongodbRepo    mongorepo.MongoDBRepo
	contextTimeout time.Duration
	redisRepo      redisrepo.RedisRepo
}

type RepoInjection struct {
	MongoDBRepo mongorepo.MongoDBRepo
	Redis       redisrepo.RedisRepo
}

func NewAppWebhookUsecase(r RepoInjection, timeout time.Duration) WebhookUsecase {
	return &webhookUsecase{
		mongodbRepo:    r.MongoDBRepo,
		contextTimeout: timeout,
		redisRepo:      r.Redis,
	}
}

type WebhookUsecase interface {
	HandleWebhook(ctx context.Context, webhookData map[string]interface{}) response.Base
}
