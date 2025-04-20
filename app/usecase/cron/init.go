package cron

import(
	"context"
	"time"
	mongorepo "app/app/repository/mongo"
	redisrepo "app/app/repository/redis"
	s3Repo "app/app/repository/s3"
)

type cronScheduler struct {
	mongodbRepo    mongorepo.MongoDBRepo
	contextTimeout time.Duration
	redisRepo      redisrepo.RedisRepo
	s3Repo         s3Repo.S3Repo
}

type RepoInjection struct {
	MongoDBRepo mongorepo.MongoDBRepo
	Redis       redisrepo.RedisRepo
	S3Repo      s3Repo.S3Repo
}

func NewCronScheduler(r RepoInjection, timeout time.Duration) CronScheduler {
    return &cronScheduler{
        mongodbRepo:    r.MongoDBRepo,
		contextTimeout: timeout,
		redisRepo:      r.Redis,
		s3Repo:         r.S3Repo,
    }
}

type CronScheduler interface {
	ConfirmNotification(ctx context.Context)
}